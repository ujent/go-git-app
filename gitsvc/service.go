package gitsvc

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ujent/go-git-app/contract"
	"github.com/ujent/go-git-mysql/mysqlfs"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/format/index"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

const filesPrefix = "files_"
const gitPrefix = "git_"

//Service - provides go-git functionality
type Service interface {
	//Filesystem returns fs of current repository
	Filesystem() (billy.Filesystem, error)

	//FilesList - returns only files in Merged mode, conflict files are excluded
	FilesList() ([]contract.FileInfo, error)

	// Repositories - returns all locally existing repositories
	Repositories() ([]string, error)

	//CreateRepository - creates a new repository
	CreateRepository(name string) error

	//OpenRepository - opens an existing repository
	OpenRepository(name string) error

	//RemoveRepository - removes specified repository permanently
	RemoveRepository(name string) error

	//CurrentRepository returns current repository name
	CurrentRepository() (name string)

	// Clone the given repository to the given directory
	Clone(url, repoName string, auth *contract.Credentials) error

	// Fetch fetches references along with the objects necessary to complete
	// their histories, from the remote named as FetchOptions.RemoteName.
	// Remote can be empty (use "origin" by default)
	//
	// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
	// no changes to be fetched, or an error.
	Fetch(remote string, auth *contract.Credentials) error

	// Pull incorporates changes from a remote repository into the current branch.
	// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
	// no changes to be fetched, or an error.
	Pull(remote string, auth *contract.Credentials) error

	// Push performs a push to the remote. Returns NoErrAlreadyUpToDate if
	// the remote was already up-to-date, from the remote named as
	// FetchOptions.RemoteName.
	// If `remote` parameter is empty, use "origin" by default
	// Use credentials if needed. Remote also can be empty
	Push(remote string, auth *contract.Credentials) error

	//Commit - commits changes and returns commit hash
	Commit(msg string) (string, error)

	//Merge - analog of git merge command
	Merge(branch string) error

	//MergeMsgShort - returns MERGE_MSG file content  with trimming strings which begin from "#"
	MergeMsgShort() (string, error)

	//MergeMsgFull - returns MERGE_MSG file content  without trimming strings which begin from "#"
	MergeMsgFull() (string, error)

	//ConflictFileList - returns pathes of files with conflicts
	ConflictFileList() ([]string, error)

	//ConflictResultFile - returns file with unresolved conflicts
	ConflictResultFile(path string) (billy.File, error)

	//FilesList - returns current repository files pathes
	ConflictFiles(path string) ([]contract.MergeFile, error)

	//Checkout - switches branch to specified commit
	Checkout(commit string) error

	//CheckoutBranch - switch to specified existing branch or creates new branch if it doesn't exist
	CheckoutBranch(branch string) error

	//CreateBranch - creates a new branch from specified commit, if commit is empty new branch will be created from current commit
	CreateBranch(branch, commit string) error

	//RemoveBranch - removes specified branch
	RemoveBranch(branch string) error

	//CurrentBranch - returns information where HEAD points now
	CurrentBranch() (*contract.Branch, error)

	//Branches - returns a list of local branches names
	Branches() ([]string, error)

	//Add - adds the file content to the staging area
	Add(path string) error

	//Log - Gets the HEAD history from HEAD, just like command "git log"
	Log() ([]contract.Commit, error)

	//CreateRemote - creates a new remote, if name isn't specified it use "origin" by default
	CreateRemote(url, name string) (*git.Remote, error)

	//RemoveRemote - delete the remote and it's config from the repository
	RemoveRemote(name string) error

	//Remotes - returns a list with all remotes
	Remotes() ([]*git.Remote, error)

	//Remote returns a remote if exists or git.ErrRemoteNotFound
	Remote(name string) (*git.Remote, error)
}

type service struct {
	user     *contract.User
	settings *contract.ServerSettings
	git      *repository
	db       *sqlx.DB
}

type repository struct {
	name string
	repo *git.Repository
	fs   billy.Filesystem
}

//New - create an instance of gitSvc
func New(user *contract.User, s *contract.ServerSettings, db *sqlx.DB) (Service, error) {

	if user.Name == "" {
		return nil, errors.New("userName cannot be empty")
	}

	if user.Email == "" {
		return nil, errors.New("userEmail cannot be empty")
	}

	return &service{user: user, settings: s, db: db}, nil
}

//Filesystem returns fs of current repository
func (svc *service) Filesystem() (billy.Filesystem, error) {
	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	return svc.git.fs, nil

}

type empty struct{}

//FilesList - returns current repository files pathes
func (svc *service) FilesList() ([]contract.FileInfo, error) {
	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	w, err := svc.git.repo.Worktree()
	if err != nil {
		return nil, err
	}

	idx, err := w.Index()
	if err != nil {
		return nil, err
	}

	res := []contract.FileInfo{}
	conf := make(map[string]struct{})

	var c empty

	for _, e := range idx.Entries {
		if e.Stage == index.Merged {
			res = append(res, contract.FileInfo{Path: e.Name})
		} else {
			conf[e.Name] = c
		}
	}

	for p := range conf {
		res = append(res, contract.FileInfo{Path: p, IsConflict: true})
	}

	return res, nil
}

//CurrentRepository returns current repository name
func (svc *service) CurrentRepository() (name string) {
	if svc.git == nil {
		return ""
	}

	return svc.git.name
}

//CreateRepository - creates a new repository
func (svc *service) CreateRepository(name string) error {

	if name == "" {
		return errors.New("Repository name cannot be empty")
	}

	fs, err := mysqlfs.New(svc.settings.GitConnStr, filesPrefix+name)
	if err != nil {
		return err
	}

	gitFs, err := mysqlfs.New(svc.settings.GitConnStr, gitPrefix+name)
	if err != nil {
		return err
	}

	st := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())
	r, err := git.Init(st, gitFs)

	svc.git = &repository{name: name, fs: gitFs, repo: r}

	if err != nil {
		return err
	}

	return nil
}

//OpenRepository - opens an existing repository
func (svc *service) OpenRepository(name string) error {
	if name == "" {
		return errors.New("Repository name cannot be empty")
	}

	if svc.git != nil && svc.git.name == name {
		return nil
	}

	fs, err := mysqlfs.New(svc.settings.GitConnStr, filesPrefix+name)
	if err != nil {
		return err
	}

	gitFs, err := mysqlfs.New(svc.settings.GitConnStr, gitPrefix+name)
	if err != nil {
		return err
	}

	st := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())

	r, err := git.Open(st, gitFs)
	if err != nil {
		return err
	}

	svc.git = &repository{name: name, fs: gitFs, repo: r}

	return nil
}

// Clone the given repository to the given directory
func (svc *service) Clone(url, repoName string, c *contract.Credentials) error {

	fs, err := mysqlfs.New(svc.settings.GitConnStr, filesPrefix+repoName)
	if err != nil {
		return err
	}

	gitFs, err := mysqlfs.New(svc.settings.GitConnStr, gitPrefix+repoName)
	if err != nil {
		return err
	}

	st := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())

	opts := &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if c != nil {
		opts.Auth = &http.BasicAuth{
			Username: c.Name,
			Password: c.Password,
		}
	}

	r, err := git.Clone(st, gitFs, opts)
	if err != nil {
		return err
	}

	svc.git = &repository{name: repoName, fs: gitFs, repo: r}

	return nil
}

// Repositories - returns all locally existing repositories
func (svc *service) Repositories() ([]string, error) {

	tables := []string{}
	svc.db.Select(&tables, "SELECT table_name FROM information_schema.tables ORDER BY table_name ASC")

	repos := []string{}
	for _, t := range tables {
		if strings.HasPrefix(t, gitPrefix) {
			repos = append(repos, strings.TrimPrefix(t, gitPrefix))
		}
	}

	return repos, nil
}

//RemoveRepository - removes specified repository permanently
func (svc *service) RemoveRepository(name string) error {

	filesTable := filesPrefix + name
	gitTable := gitPrefix + name

	tx, err := svc.db.Begin()
	if err != nil {
		return err
	}
	tx.Exec(fmt.Sprintf("DROP TABLE %s", filesTable))
	tx.Exec(fmt.Sprintf("DROP TABLE %s", gitTable))

	err = tx.Commit()
	if err != nil {
		return err
	}

	if svc.git != nil && svc.git.name == name {
		svc.git = nil
	}

	return nil
}

// Fetch fetches references along with the objects necessary to complete
// their histories, from the remote named as FetchOptions.RemoteName.
// Remote can be empty (use "origin" by default)
//
// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
// no changes to be fetched, or an error.
func (svc *service) Fetch(remote string, auth *contract.Credentials) error {

	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	if remote == "" {
		remote = "origin"
	}

	opts := &git.FetchOptions{RemoteName: remote}

	if auth != nil {
		opts.Auth = &http.BasicAuth{Username: auth.Name, Password: auth.Password}
	}

	return svc.git.repo.Fetch(opts)
}

// Pull incorporates changes from a remote repository into the current branch.
// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
// no changes to be fetched, or an error.
func (svc *service) Pull(remote string, auth *contract.Credentials) error {

	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	w, err := svc.git.repo.Worktree()
	if err != nil {
		return err
	}

	if remote == "" {
		remote = "origin"
	}

	opts := &git.PullOptions{RemoteName: remote}

	if auth != nil {
		opts.Auth = &http.BasicAuth{Username: auth.Name, Password: auth.Password}
	}

	return w.Pull(opts)
}

// Push performs a push to the remote. Returns NoErrAlreadyUpToDate if
// the remote was already up-to-date, from the remote named as
// FetchOptions.RemoteName.
// If `remote` parameter is empty, use "origin" by default
// Use credentials if needed. Remote also can be empty
func (svc *service) Push(remote string, auth *contract.Credentials) error {
	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	if remote == "" {
		remote = "origin"
	}

	opts := &git.PushOptions{RemoteName: remote}

	if auth != nil {
		opts.Auth = &http.BasicAuth{Username: auth.Name, Password: auth.Password}
	}

	return svc.git.repo.Push(opts)
}

//Commit - commits changes and returns commit hash
func (svc *service) Commit(msg string) (string, error) {
	if svc.git == nil {
		return "", contract.ErrGitRepositoryNotSet
	}

	wt, err := svc.git.repo.Worktree()
	if err != nil {
		return "", err
	}

	h, err := wt.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  svc.user.Name,
			Email: svc.user.Email,
			When:  time.Now(),
		},
	})

	if err != nil {
		return "", err
	}

	return h.String(), nil
}

//Merge - analog of git merge command
func (svc *service) Merge(branch string) error {
	if branch == "" {
		return errors.New("Branch name cannot be empty")
	}

	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	w, err := svc.git.repo.Worktree()
	if err != nil {
		return err
	}

	return w.Merge(branch)
}

//MergeMsgShort - returns MERGE_MSG file content  with trimming strings which begin from "#"
func (svc *service) MergeMsgShort() (string, error) {
	if svc.git == nil {
		return "", contract.ErrGitRepositoryNotSet
	}

	msg, err := svc.git.repo.Storer.MergeMsg()
	if err != nil {
		return "", err
	}

	return msg, nil
}

//MergeMsgFull - returns MERGE_MSG file content  without trimming strings which begin from "#"
func (svc *service) MergeMsgFull() (string, error) {
	if svc.git == nil {
		return "", contract.ErrGitRepositoryNotSet
	}

	msg, err := svc.git.repo.Storer.MergeMsgFileContent()
	if err != nil {
		return "", err
	}

	return msg, nil
}

//ConflictFileList - returns pathes of files with conflicts
func (svc *service) ConflictFileList() ([]string, error) {
	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	w, err := svc.git.repo.Worktree()
	if err != nil {
		return nil, err
	}

	withConflicts, err := w.ConflictEntries()
	if err != nil {
		return nil, err
	}

	pathes := []string{}
	for path := range withConflicts {
		pathes = append(pathes, path)
	}

	return pathes, nil
}

//ConflictResultFile - returns file with unresolved conflicts
func (svc *service) ConflictResultFile(path string) (billy.File, error) {
	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	w, err := svc.git.repo.Worktree()
	if err != nil {
		return nil, err
	}

	f, err := w.Filesystem.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}

//ConflictFiles - returns base, ours or theirs files by path of conflict file
func (svc *service) ConflictFiles(path string) ([]contract.MergeFile, error) {
	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	w, err := svc.git.repo.Worktree()
	if err != nil {
		return nil, err
	}

	withConflicts, err := w.ConflictEntries()
	if err != nil {
		return nil, err
	}

	res := []contract.MergeFile{}
	for p, entries := range withConflicts {

		if path == p {
			for _, e := range entries {
				//read base, ours or theirs files
				f, err := w.ReadFileByStage(path, e.Stage)
				if err != nil {
					return nil, err
				}

				res = append(res, contract.MergeFile{Path: p, Stage: svc.toFileStage(e.Stage), Reader: f})
			}

			break
		}

	}

	return res, nil
}

func (svc *service) toFileStage(st index.Stage) contract.FileStage {
	switch st {
	case index.Merged:
		return contract.Merged
	case index.AncestorMode:
		return contract.AncestorMode
	case index.OurMode:
		return contract.OurMode
	case index.TheirMode:
		return contract.TheirMode
	default:
		return contract.Unexpected
	}
}

//Checkout - switches branch to specified commit
func (svc *service) Checkout(commit string) error {
	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	wt, err := svc.git.repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	})

	if err != nil {
		return err
	}

	return nil
}

//CheckoutBranch - switch to specified existing branch or creates new branch if it doesn't exist
func (svc *service) CheckoutBranch(branch string) error {
	if branch == "" {
		return errors.New("Branch name cannot be empty")
	}

	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	iter, err := svc.git.repo.Branches()
	if err != nil {
		return err
	}

	var hasBranch bool

	err = iter.ForEach(func(br *plumbing.Reference) error {
		if br.Name().Short() == branch {
			hasBranch = true
			iter.Close()
		}

		return nil
	})

	if err != nil {
		return err
	}

	if !hasBranch {
		return svc.CreateBranch(branch, "")
	}

	wt, err := svc.git.repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})

	if err != nil {
		return err
	}

	return nil
}

//CreateBranch - creates a new branch from specified commit, if commit is empty new branch will be created from current commit
func (svc *service) CreateBranch(branch, commit string) error {
	if branch == "" {
		return errors.New("Branch name cannot be empty")
	}

	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	wt, err := svc.git.repo.Worktree()
	if err != nil {
		return err
	}

	var hash plumbing.Hash

	if commit == "" {

		headRef, err := svc.git.repo.Head()
		if err != nil {
			return err
		}

		hash = headRef.Hash()
	} else {
		hash = plumbing.NewHash(commit)
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Hash:   hash,
		Branch: plumbing.NewBranchReferenceName(branch),
		Create: true,
	})

	if err != nil {
		return err
	}

	return nil
}

//RemoveBranch - removes specified branch
func (svc *service) RemoveBranch(branch string) error {
	if branch == "" {
		return errors.New("Branch name cannot be empty")
	}

	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	ref := plumbing.NewBranchReferenceName(branch)

	return svc.git.repo.Storer.RemoveReference(ref)
}

//CurrentBranch - returns information where HEAD points now
func (svc *service) CurrentBranch() (*contract.Branch, error) {
	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	headRef, err := svc.git.repo.Head()
	if err != nil {
		return nil, err
	}

	return &contract.Branch{Name: headRef.Name().Short(), Hash: headRef.Hash().String()}, nil
}

//Branches - returns a list of local branches names
func (svc *service) Branches() ([]string, error) {
	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	iter, err := svc.git.repo.Branches()
	if err != nil {
		return nil, err
	}

	branches := []string{}

	err = iter.ForEach(func(br *plumbing.Reference) error {
		if !br.Name().IsRemote() {
			branches = append(branches, br.Name().Short())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return branches, nil
}

//Add - adds the file content to the staging area
func (svc *service) Add(path string) error {

	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	wt, err := svc.git.repo.Worktree()
	if err != nil {
		return nil
	}

	return wt.Add(path)
}

//Log - Gets the HEAD history from HEAD, just like command "git log"
func (svc *service) Log() ([]contract.Commit, error) {

	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	ref, err := svc.git.repo.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := svc.git.repo.Log(&git.LogOptions{From: ref.Hash()})
	res := []contract.Commit{}

	err = cIter.ForEach(func(c *object.Commit) error {

		res = append(res, contract.Commit{Author: &contract.User{Name: c.Author.Name, Email: c.Author.Email}, Date: c.Author.When, Message: c.Message, Hash: c.Hash.String()})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

//CreateRemote - creates a new remote, if name isn't specified it use "origin" by default
func (svc *service) CreateRemote(url, name string) (*git.Remote, error) {
	if url == "" {
		return nil, errors.New("Remote url cannot be empty")
	}

	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	if name == "" {
		name = "origin"
	}

	r, err := svc.git.repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})

	if err != nil {
		return nil, err
	}

	return r, nil
}

//RemoveRemote - delete the remote and it's config from the repository
func (svc *service) RemoveRemote(name string) error {
	if name == "" {
		return errors.New("Remote name cannot be empty")
	}

	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	return svc.git.repo.DeleteRemote(name)
}

//Remotes - returns a list with all remotes
func (svc *service) Remotes() ([]*git.Remote, error) {
	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	remotes, err := svc.git.repo.Remotes()
	if err != nil {
		return nil, err
	}

	return remotes, nil
}

//Remote returns a remote if exists or git.ErrRemoteNotFound
func (svc *service) Remote(name string) (*git.Remote, error) {
	if name == "" {
		return nil, errors.New("Remote name cannot be empty")
	}

	if svc.git == nil {
		return nil, contract.ErrGitRepositoryNotSet
	}

	r, err := svc.git.repo.Remote(name)
	if err != nil {
		return nil, err
	}

	return r, nil
}
