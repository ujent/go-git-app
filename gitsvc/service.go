package gitsvc

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	Repositories() ([]string, error)
	CreateRepository(name string) error
	OpenRepository(name string) error
	RemoveRepository(name string) error
	CurrentRepository() (name string)
	Clone(url, repoName string, auth *contract.Credentials) error
	Fetch(remote string) error
	Pull() error
	Push(remote string, auth *contract.Credentials) error
	Commit(msg string) (string, error)
	Merge(branch string) error
	MergeMsgShort() (string, error)
	MergeMsgFull() (string, error)
	ConflictFileList() ([]string, error)
	ConflictResultFile(path string) (billy.File, error)
	ConflictFiles(path string) ([]contract.MergeFile, error)
	Checkout(commit string) error
	CheckoutBranch(branch string) error
	CreateBranch(branch, commit string) error
	RemoveBranch(branch string) error
	CurrentBranch() (*contract.Branch, error)
	Branches() ([]string, error)
	Add(path string) error
	Log() ([]contract.Commit, error)
	Filesystem() billy.Filesystem
	CreateRemote(url, name string) error
	RemoveRemote(name string) error
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

func (svc *service) Filesystem() billy.Filesystem {
	if svc.git != nil {
		return svc.git.fs
	}

	return nil
}

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
func (svc *service) Fetch(remote string) error {

	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	err := svc.git.repo.Fetch(&git.FetchOptions{
		RemoteName: remote,
	})

	if err != nil {
		return err
	}

	return nil
}

// Pull incorporates changes from a remote repository into the current branch.
// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
// no changes to be fetched, or an error.
func (svc *service) Pull() error {

	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	w, err := svc.git.repo.Worktree()
	if err != nil {
		return err
	}

	return w.Pull(&git.PullOptions{RemoteName: "origin"})
}

// Push performs a push to the remote. Returns NoErrAlreadyUpToDate if
// the remote was already up-to-date, from the remote named as
// FetchOptions.RemoteName.
// Use credentials if needed. Remote also can be empty
func (svc *service) Push(remote string, auth *contract.Credentials) error {
	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
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

				bytes, err := ioutil.ReadAll(f)
				if err != nil && err != io.EOF {
					return nil, err
				}

				res = append(res, contract.MergeFile{Path: p, Stage: svc.toFileStage(e.Stage), Content: string(bytes)})
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
func (svc *service) CreateRemote(url, name string) error {
	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	if url == "" {
		return errors.New("Remote url cannot be empty")
	}

	if name == "" {
		name = "origin"
	}

	_, err := svc.git.repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})

	if err != nil {
		return err
	}

	return nil
}

//RemoveRemote - delete the remote and it's config from the repository
func (svc *service) RemoveRemote(name string) error {
	if svc.git == nil {
		return contract.ErrGitRepositoryNotSet
	}

	if name == "" {
		return errors.New("Remote name cannot be empty")
	}

	return svc.git.repo.DeleteRemote(name)
}
