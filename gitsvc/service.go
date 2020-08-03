package gitsvc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	git "bitbucket.org/vishjosh/bipp-go-git"
	"bitbucket.org/vishjosh/bipp-go-git/config"
	"bitbucket.org/vishjosh/bipp-go-git/experimental-app/contract"
	"bitbucket.org/vishjosh/bipp-go-git/go-billy"
	"bitbucket.org/vishjosh/bipp-go-git/go-billy/osfs"
	"bitbucket.org/vishjosh/bipp-go-git/plumbing"
	"bitbucket.org/vishjosh/bipp-go-git/plumbing/cache"
	"bitbucket.org/vishjosh/bipp-go-git/plumbing/format/index"
	"bitbucket.org/vishjosh/bipp-go-git/plumbing/object"
	"bitbucket.org/vishjosh/bipp-go-git/plumbing/transport/http"
	"bitbucket.org/vishjosh/bipp-go-git/storage/filesystem"
	"bitbucket.org/vishjosh/bipp-go-git/storage/mysqlfs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const filesPrefix = "files_"
const gitPrefix = "git_"

//Service - provides go-git functionality
type Service interface {

	//CurrentUser - current login user
	CurrentUser() *contract.User

	//SwitchUser - change user from which we are using app
	SwitchUser(user *contract.User) error

	//Filesystem returns fs of current repository
	Filesystem(user, repo string) (billy.Filesystem, error)

	//FilesList - returns current repository files pathes
	FilesList(rq *contract.BaseRequest) ([]contract.FileInfo, error)

	// Repositories - returns all locally existing repositories
	Repositories(user string) ([]string, error)

	//CreateRepository - creates a new repository
	CreateRepository(user, repo string) error

	//OpenRepository - opens an existing repository
	OpenRepository(user, repo string) error

	//RemoveRepository - removes specified repository permanently
	RemoveRepository(user, repo string) error

	//CurrentRepository returns current repository name
	CurrentRepository() (name string)

	// Clone the given repository to the given directory
	Clone(user, url string, auth *contract.Credentials) (string, error)

	// Fetch fetches references along with the objects necessary to complete
	// their histories, from the remote named as FetchOptions.RemoteName.
	// Remote can be empty (use "origin" by default)
	//
	// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
	// no changes to be fetched, or an error.
	Fetch(user, repo, remote string, auth *contract.Credentials) error

	// Pull incorporates changes from a remote repository into the current branch.
	// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
	// no changes to be fetched, or an error.
	Pull(rq *contract.BaseRequest, remote string, auth *contract.Credentials) (string, error)

	// Push performs a push to the remote. Returns NoErrAlreadyUpToDate if
	// the remote was already up-to-date, from the remote named as
	// FetchOptions.RemoteName.
	// If `remote` parameter is empty, use "origin" by default
	// Use credentials if needed. Remote also can be empty
	Push(rq *contract.BaseRequest, remote string, auth *contract.Credentials) error

	//Commit - commits changes and returns commit hash
	Commit(rq *contract.BaseRequest, msg string) (string, error)

	//Merge - analog of git merge command
	Merge(rq *contract.BaseRequest, branch string) (string, error)

	//MergeMsgShort - returns MERGE_MSG file content  with trimming strings which begin from "#"
	MergeMsgShort(rq *contract.BaseRequest) (string, error)

	//MergeMsgFull - returns MERGE_MSG file content  without trimming strings which begin from "#"
	MergeMsgFull(rq *contract.BaseRequest) (string, error)

	//AbortMerge will abort the merge process and try to reconstruct the pre-merge state
	AbortMerge(rq *contract.BaseRequest) error

	//ConflictFileList - returns pathes of files with conflicts
	ConflictFileList(rq *contract.BaseRequest) ([]string, error)

	//FilesList - returns current repository files pathes
	ConflictFiles(rq *contract.BaseRequest, path string) ([]contract.MergeFile, error)

	//Checkout - switches branch to specified commit
	Checkout(user, repo string, commit string) error

	//CheckoutBranch - switch to specified existing branch or creates new branch if it doesn't exist
	CheckoutBranch(user, repo, branch string) error

	//CreateBranch - creates a new branch from specified commit, if commit is empty new branch will be created from current commit
	CreateBranch(user, repo, branch, commit string) error

	//RemoveBranch - removes specified branch
	RemoveBranch(user, repo, branch string) error

	//CurrentBranch - returns information where HEAD points now
	CurrentBranch() (*contract.Branch, error)

	//Branches - returns a list of local branches names
	Branches(user, repo string) ([]string, error)

	//Add - adds the file content to the staging area
	Add(rq *contract.BaseRequest, path string) error

	//Log - Gets the HEAD history from HEAD, just like command "git log"
	Log(rq *contract.BaseRequest) ([]contract.Commit, error)

	//CreateRemote - creates a new remote, if name isn't specified it use "origin" by default
	CreateRemote(user, repo, url, name string) (*git.Remote, error)

	//RemoveRemote - delete the remote and it's config from the repository
	RemoveRemote(user, repo, name string) error

	//Remotes - returns a list with all remotes
	Remotes(user, repo string) ([]*git.Remote, error)

	//Remote returns a remote if exists or git.ErrRemoteNotFound
	Remote(user, repo, name string) (*git.Remote, error)

	//File - returns billy file
	File(rq *contract.BaseRequest, path string) (billy.File, error)

	//AddFile - add file to Filesystem and stage it
	AddFile(rq *contract.BaseRequest, path, content string) error

	//EditFile - save changed in file and stage changes
	EditFile(rq *contract.BaseRequest, path, content string) error

	//RemoveFile - remove file from Filesystem and stage changes
	RemoveFile(rq *contract.BaseRequest, path string) error

	// Status - returns the working tree status
	Status(rq *contract.BaseRequest) (git.Status, error)
}

type service struct {
	user     *contract.User
	settings *contract.ServerSettings
	git      *repository
	db       *sqlx.DB
}

type repository struct {
	sync.Mutex
	name string
	repo *git.Repository
	fs   billy.Filesystem
}

//New - create an instance of gitSvc
func New(s *contract.ServerSettings, db *sqlx.DB) (Service, error) {

	return &service{user: &contract.User{}, settings: s, db: db}, nil
}

func (svc *service) File(rq *contract.BaseRequest, path string) (billy.File, error) {
	err := svc.validateBaseRQ(rq)
	if err != nil {
		return nil, err
	}

	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	err = svc.setSettings(&contract.User{Name: rq.User.Name}, rq.Repository, rq.Branch)
	if err != nil {
		return nil, err
	}

	fs := svc.git.fs
	f, err := fs.OpenFile(path, os.O_RDONLY, 0666)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (svc *service) AddFile(rq *contract.BaseRequest, path, content string) error {
	if rq == nil {
		return errors.New("rq cannot be nil")
	}

	if rq.User == nil {
		return errors.New("User cannot be empty")
	}

	if rq.Repository == "" {
		return errors.New("Repository cannot be empty")
	}

	if path == "" {
		return errors.New("path cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: rq.User.Name}, rq.Repository, rq.Branch)
	if err != nil {
		return err
	}

	svc.git.Lock()
	defer svc.git.Unlock()

	fs := svc.git.fs
	_, err = fs.Create(path)

	if err != nil {
		return err
	}

	wt, err := svc.git.repo.Worktree()
	if err != nil {
		return nil
	}

	return wt.Add(path)
}

func (svc *service) EditFile(rq *contract.BaseRequest, path, content string) error {
	err := svc.validateBaseRQ(rq)
	if err != nil {
		return err
	}

	if path == "" {
		return errors.New("path cannot be empty")
	}

	err = svc.setSettings(&contract.User{Name: rq.User.Name}, rq.Repository, rq.Branch)
	if err != nil {
		return err
	}

	var f billy.File

	fs := svc.git.fs
	f, err = fs.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0666)

	if err != nil {
		return err
	}

	_, err = f.Write([]byte(content))
	if err != nil {
		return err
	}

	f.Close()

	wt, err := svc.git.repo.Worktree()
	if err != nil {
		return nil
	}

	return wt.Add(path)
}

func (svc *service) RemoveFile(rq *contract.BaseRequest, path string) error {
	err := svc.validateBaseRQ(rq)
	if err != nil {
		return err
	}

	if path == "" {
		return errors.New("path cannot be empty")
	}

	err = svc.setSettings(&contract.User{Name: rq.User.Name}, rq.Repository, rq.Branch)
	if err != nil {
		return err
	}

	wt, err := svc.git.repo.Worktree()
	if err != nil {
		return err
	}

	fs := svc.git.fs
	err = fs.Remove(path)

	if err != nil {
		return err
	}

	return wt.Add(path)
}

func (svc *service) Status(rq *contract.BaseRequest) (git.Status, error) {
	err := svc.validateBaseRQ(rq)
	if err != nil {
		return nil, err
	}

	err = svc.setSettings(&contract.User{Name: rq.User.Name}, rq.Repository, rq.Branch)
	if err != nil {
		return nil, err
	}

	wt, err := svc.git.repo.Worktree()
	if err != nil {
		return nil, err
	}

	return wt.Status()
}

func (svc *service) SwitchUser(user *contract.User) error {

	if user == nil {
		return errors.New("user cannot be empty")
	}

	if user.Name == "" {
		return errors.New("userName cannot be empty")
	}

	if user.Email == "" {
		user.Email = user.Name + "@test.com"
	}

	if user.Name == svc.user.Name {
		svc.user.Email = user.Email

		return nil
	}

	svc.user = user
	svc.git = nil

	return nil
}

func (svc *service) CurrentUser() *contract.User {
	return svc.user
}

func (svc *service) setSettings(user *contract.User, repo, branch string) error {

	err := svc.SwitchUser(user)
	if err != nil {
		return err
	}

	if svc.git == nil || (svc.git != nil && svc.git.name != repo) {

		err = svc.OpenRepository(user.Name, repo)
		if err != nil {
			return err
		}
	}

	if branch != "" {
		currBr, err := svc.CurrentBranch()
		if err != nil {
			return err
		}

		if currBr.Name != branch {
			err = svc.CheckoutBranch(user.Name, repo, branch)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (svc *service) validateBaseRQ(rq *contract.BaseRequest) error {
	if rq == nil {
		return errors.New("rq cannot be nil")
	}

	if rq.User == nil {
		return errors.New("User cannot be empty")
	}

	if rq.Repository == "" {
		return errors.New("Repository cannot be empty")
	}

	if rq.Branch == "" {
		return errors.New("Branch cannot be empty")
	}

	return nil
}

func (svc *service) validateBaseRQWithoutBranch(rq *contract.BaseRequest) error {
	if rq == nil {
		return errors.New("rq cannot be nil")
	}

	if rq.User == nil {
		return errors.New("User cannot be empty")
	}

	if rq.Repository == "" {
		return errors.New("Repository cannot be empty")
	}

	return nil
}

//Filesystem returns fs of current repository
func (svc *service) Filesystem(user, repo string) (billy.Filesystem, error) {
	if user == "" {
		return nil, errors.New("User cannot be empty")
	}

	if repo == "" {
		return nil, errors.New("Repository cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: user}, repo, "")
	if err != nil {
		return nil, err
	}

	return svc.git.fs, nil

}

type empty struct{}

//FilesList - returns current repository files pathes
func (svc *service) FilesList(rq *contract.BaseRequest) ([]contract.FileInfo, error) {

	err := svc.validateBaseRQ(rq)
	if err != nil {
		return nil, err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return nil, err
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
func (svc *service) CreateRepository(user, repo string) error {

	if repo == "" {
		return errors.New("Repository name cannot be empty")
	}

	err := svc.SwitchUser(&contract.User{Name: user})
	if err != nil {
		return err
	}

	fs, gitFs, err := svc.createFs(user, repo)
	if err != nil {
		return err
	}

	st := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())
	r, err := git.Init(st, gitFs)

	svc.git = &repository{name: repo, fs: gitFs, repo: r}

	if err != nil {
		return err
	}

	return nil
}

func (svc *service) tablesNames(user, repoName string) (filesTableName, gitTableName string, err error) {
	if svc.user == nil {
		return "", "", errors.New("user cannot be empty")
	}

	if user == "" {
		return "", "", errors.New("userName cannot be empty")
	}

	filesTableName = filesPrefix + user + "_" + repoName
	gitTableName = gitPrefix + user + "_" + repoName

	return filesTableName, gitTableName, nil
}

func (svc *service) gitPath(user, repoName string) (gitPath, wtPath string, err error) {
	if svc.user == nil {
		return "", "", errors.New("user cannot be empty")
	}

	if user == "" {
		return "", "", errors.New("userName cannot be empty")
	}

	wtPath = filepath.Join(svc.settings.GitRoot, user, repoName)
	gitPath = filepath.Join(wtPath, ".git")

	return gitPath, wtPath, nil
}

//OpenRepository - opens an existing repository
func (svc *service) OpenRepository(user, repo string) error {

	if repo == "" {
		return errors.New("Repository name cannot be empty")
	}

	err := svc.SwitchUser(&contract.User{Name: user})
	if err != nil {
		return err
	}

	if svc.git != nil && svc.git.name == repo {
		return nil
	}

	fs, gitFs, err := svc.createFs(user, repo)
	if err != nil {
		return err
	}

	st := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())

	r, err := git.Open(st, gitFs)
	if err != nil {
		return err
	}

	svc.git = &repository{name: repo, fs: gitFs, repo: r}

	return nil
}

func (svc *service) createFs(user, repo string) (fs billy.Filesystem, gitFs billy.Filesystem, err error) {
	switch svc.settings.FsType {
	case contract.FsTypeMySQL:
		{
			filesTableName, gitTableName, err := svc.tablesNames(user, repo)
			if err != nil {
				return nil, nil, err
			}

			fs, err = mysqlfs.New(svc.db.DB, gitTableName)
			if err != nil {
				return nil, nil, err
			}

			gitFs, err = mysqlfs.New(svc.db.DB, filesTableName)
			if err != nil {
				return nil, nil, err
			}

			return fs, gitFs, nil
		}
	case contract.FsTypeLocal:
		{
			filesPath, wtPath, err := svc.gitPath(user, repo)
			if err != nil {
				return nil, nil, err
			}

			dotgit := osfs.New(filesPath)
			wtfs := osfs.New(wtPath)

			return dotgit, wtfs, nil
		}
	default:
		{
			return nil, nil, fmt.Errorf("Wrong fsType = %d", svc.settings.FsType)
		}
	}
}

// Clone the given repository to the given directory
func (svc *service) Clone(user, url string, auth *contract.Credentials) (string, error) {

	if url == "" {
		return "", errors.New("URL cannot be empty")
	}

	splitted := strings.Split(url, "/")
	last := splitted[len(splitted)-1]
	repoName := strings.TrimSuffix(last, ".git")

	if repoName == "" {
		return "", errors.New("wrong URL for clone operation")
	}

	err := svc.SwitchUser(&contract.User{Name: user})
	if err != nil {
		return "", err
	}

	fs, gitFs, err := svc.createFs(user, repoName)
	if err != nil {
		return "", err
	}

	st := filesystem.NewStorage(fs, cache.NewObjectLRUDefault())

	opts := &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if auth != nil {
		opts.Auth = &http.BasicAuth{
			Username: auth.Name,
			Password: auth.Password,
		}
	}

	r, err := git.Clone(st, gitFs, opts)

	if err != nil {
		delErr := svc.deleteRepo(user, repoName)
		if delErr != nil {
			log.Printf("Cannot remove repo: %s, user: %s, error: %v\n", repoName, user, err)
		}

		return "", err
	}

	svc.git = &repository{name: repoName, fs: gitFs, repo: r}

	return repoName, nil
}

func (svc *service) deleteRepo(user, repo string) error {
	switch svc.settings.FsType {
	case contract.FsTypeMySQL:
		{
			filesTable, gitTable, err := svc.tablesNames(user, repo)
			if err != nil {
				return err
			}

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
		}
	case contract.FsTypeLocal:
		{
			p := filepath.Join(svc.settings.GitRoot, user, repo)

			err := removeFiles(p)
			if err != nil {
				return err
			}
		}
	default:
		{
			return fmt.Errorf("[Delete Repo] Wrong fsType = %d", svc.settings.FsType)
		}
	}

	return nil
}

func removeFiles(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}

	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return nil
}

// Repositories - returns all locally existing repositories
func (svc *service) Repositories(user string) ([]string, error) {
	if user == "" {
		return nil, errors.New("user cannot be empty")
	}

	switch svc.settings.FsType {
	case contract.FsTypeMySQL:
		{

			tables := []string{}
			svc.db.Select(&tables, "SELECT table_name FROM information_schema.tables ORDER BY table_name ASC")

			repos := []string{}

			for _, t := range tables {
				if strings.HasPrefix(t, gitPrefix) {
					temp := strings.TrimPrefix(t, gitPrefix)
					userPrefix := user + "_"

					if strings.HasPrefix(temp, userPrefix) {
						repos = append(repos, strings.TrimPrefix(temp, userPrefix))
					}
				}
			}

			return repos, nil
		}
	case contract.FsTypeLocal:
		{
			gitPath := filepath.Join(svc.settings.GitRoot, user)
			files, err := ioutil.ReadDir(gitPath)

			if err != nil {
				if os.IsNotExist(err) {
					return []string{}, nil
				}

				return nil, err
			}

			res := []string{}

			for _, f := range files {
				if f.IsDir() {
					res = append(res, f.Name())
				}
			}

			return res, nil
		}
	default:
		{
			return nil, fmt.Errorf("Wrong fsType = %d", svc.settings.FsType)
		}
	}

}

//RemoveRepository - removes specified repository permanently
func (svc *service) RemoveRepository(user, repo string) error {
	err := svc.SwitchUser(&contract.User{Name: user})
	if err != nil {
		return err
	}

	err = svc.deleteRepo(user, repo)
	if err != nil {
		return err
	}

	if svc.git != nil && svc.git.name == repo {
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
func (svc *service) Fetch(user, repo, remote string, auth *contract.Credentials) error {

	err := svc.setSettings(&contract.User{Name: user}, repo, "")
	if err != nil {
		return err
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
func (svc *service) Pull(rq *contract.BaseRequest, remote string, auth *contract.Credentials) (string, error) {

	err := svc.validateBaseRQWithoutBranch(rq)
	if err != nil {
		return "", err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return "", err
	}

	w, err := svc.git.repo.Worktree()
	if err != nil {
		return "", err
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
func (svc *service) Push(rq *contract.BaseRequest, remote string, auth *contract.Credentials) error {

	err := svc.validateBaseRQWithoutBranch(rq)
	if err != nil {
		return err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return err
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
func (svc *service) Commit(rq *contract.BaseRequest, msg string) (string, error) {

	err := svc.validateBaseRQWithoutBranch(rq)
	if err != nil {
		return "", err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return "", err
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
func (svc *service) Merge(rq *contract.BaseRequest, branch string) (string, error) {
	if branch == "" {
		return "", errors.New("Branch name cannot be empty")
	}

	err := svc.validateBaseRQ(rq)
	if err != nil {
		return "", err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return "", err
	}

	w, err := svc.git.repo.Worktree()
	if err != nil {
		return "", err
	}

	return w.Merge(branch)
}

//AbortMerge will abort the merge process and try to reconstruct the pre-merge state
func (svc *service) AbortMerge(rq *contract.BaseRequest) error {
	err := svc.validateBaseRQ(rq)
	if err != nil {
		return err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return err
	}

	w, err := svc.git.repo.Worktree()
	if err != nil {
		return err
	}

	return w.AbortMerge()
}

//MergeMsgShort - returns MERGE_MSG file content  with trimming strings which begin from "#"
func (svc *service) MergeMsgShort(rq *contract.BaseRequest) (string, error) {
	err := svc.validateBaseRQ(rq)
	if err != nil {
		return "", err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return "", err
	}

	msg, err := svc.git.repo.Storer.MergeMsg()
	if err != nil {
		return "", err
	}

	return msg, nil
}

//MergeMsgFull - returns MERGE_MSG file content  without trimming strings which begin from "#"
func (svc *service) MergeMsgFull(rq *contract.BaseRequest) (string, error) {
	err := svc.validateBaseRQ(rq)
	if err != nil {
		return "", err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return "", err
	}

	msg, err := svc.git.repo.Storer.MergeMsgFileContent()
	if err != nil {
		return "", err
	}

	return msg, nil
}

//ConflictFileList - returns pathes of files with conflicts
func (svc *service) ConflictFileList(rq *contract.BaseRequest) ([]string, error) {
	err := svc.validateBaseRQ(rq)
	if err != nil {
		return nil, err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return nil, err
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

//ConflictFiles - returns base, ours or theirs files by path of conflict file
func (svc *service) ConflictFiles(rq *contract.BaseRequest, path string) ([]contract.MergeFile, error) {
	err := svc.validateBaseRQ(rq)
	if err != nil {
		return nil, err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return nil, err
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
func (svc *service) Checkout(user, repo string, commit string) error {
	if user == "" {
		return errors.New("User cannot be empty")
	}

	if repo == "" {
		return errors.New("Repository cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: user}, repo, "")
	if err != nil {
		return err
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
func (svc *service) CheckoutBranch(user, repo, branch string) error {
	if branch == "" {
		return errors.New("Branch name cannot be empty")
	}

	if user == "" {
		return errors.New("User cannot be empty")
	}

	if repo == "" {
		return errors.New("Repository cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: user}, repo, "")
	if err != nil {
		return err
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
		return svc.CreateBranch(user, repo, branch, "")
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
func (svc *service) CreateBranch(user, repo, branch, commit string) error {
	if branch == "" {
		return errors.New("Branch name cannot be empty")
	}

	if user == "" {
		return errors.New("User cannot be empty")
	}

	if repo == "" {
		return errors.New("Repository cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: user}, repo, "")
	if err != nil {
		return err
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
func (svc *service) RemoveBranch(user, repo, branch string) error {
	if branch == "" {
		return errors.New("Branch name cannot be empty")
	}

	if user == "" {
		return errors.New("User cannot be empty")
	}

	if repo == "" {
		return errors.New("Repository cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: user}, repo, "")
	if err != nil {
		return err
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
		if err == plumbing.ErrReferenceNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &contract.Branch{Name: headRef.Name().Short(), Hash: headRef.Hash().String()}, nil
}

//Branches - returns a list of local branches names
func (svc *service) Branches(user, repo string) ([]string, error) {

	if user == "" {
		return nil, errors.New("User cannot be empty")
	}

	if repo == "" {
		return nil, errors.New("Repository cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: user}, repo, "")

	if err != nil {
		return nil, err
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
func (svc *service) Add(rq *contract.BaseRequest, path string) error {

	err := svc.validateBaseRQWithoutBranch(rq)
	if err != nil {
		return err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return nil
	}

	wt, err := svc.git.repo.Worktree()
	if err != nil {
		return nil
	}

	return wt.Add(path)
}

//Log - Gets the HEAD history from HEAD, just like command "git log"
func (svc *service) Log(rq *contract.BaseRequest) ([]contract.Commit, error) {

	err := svc.validateBaseRQ(rq)
	if err != nil {
		return nil, err
	}

	err = svc.setSettings(rq.User, rq.Repository, rq.Branch)
	if err != nil {
		return nil, nil
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
func (svc *service) CreateRemote(user, repo, url, name string) (*git.Remote, error) {
	if url == "" {
		return nil, errors.New("Remote url cannot be empty")
	}

	if user == "" {
		return nil, errors.New("User cannot be empty")
	}

	if repo == "" {
		return nil, errors.New("Repository cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: user}, repo, "")
	if err != nil {
		return nil, err
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
func (svc *service) RemoveRemote(user, repo, name string) error {
	if name == "" {
		return errors.New("Remote name cannot be empty")
	}

	if user == "" {
		return errors.New("User cannot be empty")
	}

	if repo == "" {
		return errors.New("Repository cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: user}, repo, "")
	if err != nil {
		return err
	}

	return svc.git.repo.DeleteRemote(name)
}

//Remotes - returns a list with all remotes
func (svc *service) Remotes(user, repo string) ([]*git.Remote, error) {
	if user == "" {
		return nil, errors.New("User cannot be empty")
	}

	if repo == "" {
		return nil, errors.New("Repository cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: user}, repo, "")
	if err != nil {
		return nil, err
	}

	remotes, err := svc.git.repo.Remotes()
	if err != nil {
		return nil, err
	}

	return remotes, nil
}

//Remote returns a remote if exists or git.ErrRemoteNotFound
func (svc *service) Remote(user, repo, name string) (*git.Remote, error) {
	if name == "" {
		return nil, errors.New("Remote name cannot be empty")
	}

	if user == "" {
		return nil, errors.New("User cannot be empty")
	}

	if repo == "" {
		return nil, errors.New("Repository cannot be empty")
	}

	err := svc.setSettings(&contract.User{Name: user}, repo, "")
	if err != nil {
		return nil, err
	}

	r, err := svc.git.repo.Remote(name)
	if err != nil {
		return nil, err
	}

	return r, nil
}
