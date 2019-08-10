package gitsvc

import (
	"errors"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ujent/go-git-app/contract"
	"github.com/ujent/go-git-mysql/mysqlfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

const filesPrefix = "files_"
const gitPrefix = "git_"

//Service - provides go-git functionality
type Service interface {
	//Remove repository is not supported by go-git

	CreateRepository(name string) error
	OpenRepository(name string) error
	Repositories() ([]string, error)
	Clone(url, repoName string, auth *contract.Credentials) error
	Fetch(remote string) error
	Pull() error
	Push(remote string, auth *contract.Credentials) error
	Commit(msg string) error
	Merge(branch string) error
	Branches() ([]string, error)
	Checkout(commit string) error
	CheckoutBranch(branch string) error
	CreateBranch(branch, commit string) error
	//CreateRemoteBranch(branch string) error
	DeleteBranch(branch string) error
	Add(path string) error
	Log() ([]contract.Commit, error)
}

type service struct {
	user     *contract.User
	settings *contract.ServerSettings
	gitRepo  *git.Repository
	db       *sqlx.DB
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
	svc.gitRepo = r

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

	svc.gitRepo = r

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

	svc.gitRepo = r

	return nil
}

// Repositories - returns all locally existing repositories
func (svc *service) Repositories() ([]string, error) {

	tables := []string{}
	svc.db.Select(&tables, "SELECT table_name FROM information_schema.tables WHERE table_type = 'base table' ORDER BY table_name ASC")

	repos := []string{}
	for _, t := range tables {
		if strings.HasPrefix(t, gitPrefix) {
			repos = append(repos, strings.TrimPrefix(t, gitPrefix))
		}
	}

	return repos, nil
}

//Branches - returns a list of local branches names
func (svc *service) Branches() ([]string, error) {
	iter, err := svc.gitRepo.Branches()
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

// Fetch fetches references along with the objects necessary to complete
// their histories, from the remote named as FetchOptions.RemoteName.
// Remote can be empty (use "origin" by default)
//
// Returns nil if the operation is successful, NoErrAlreadyUpToDate if there are
// no changes to be fetched, or an error.
func (svc *service) Fetch(remote string) error {

	err := svc.gitRepo.Fetch(&git.FetchOptions{
		RemoteName: remote,
	})

	if err != nil {
		return err
	}

	return nil
}

func (svc *service) Pull() error {
	return nil
}

//Push performs a push to the remote. Returns NoErrAlreadyUpToDate if
// the remote was already up-to-date, from the remote named as
// FetchOptions.RemoteName.
//Use credentials if needed. Remote also can be empty
func (svc *service) Push(remote string, auth *contract.Credentials) error {
	opts := &git.PushOptions{RemoteName: remote}

	if auth != nil {
		opts.Auth = &http.BasicAuth{Username: auth.Name, Password: auth.Password}
	}

	return svc.gitRepo.Push(opts)
}

//Commit - commits changes and log history
func (svc *service) Commit(msg string) error {
	wt, err := svc.gitRepo.Worktree()
	if err != nil {
		return err
	}

	_, err = wt.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  svc.user.Email,
			Email: svc.user.Email,
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (svc *service) Merge(branch string) error {
	return nil
}

//Checkout - switches branch to specified commit
func (svc *service) Checkout(commit string) error {
	wt, err := svc.gitRepo.Worktree()
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

//CheckoutBranch - switch to specified existing branch
func (svc *service) CheckoutBranch(branch string) error {
	wt, err := svc.gitRepo.Worktree()
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
	wt, err := svc.gitRepo.Worktree()
	if err != nil {
		return err
	}

	var hash plumbing.Hash

	if commit == "" {

		headRef, err := svc.gitRepo.Head()
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

//DeleteBranch - removes specified branch
func (svc *service) DeleteBranch(branch string) error {
	ref := plumbing.NewBranchReferenceName(branch)

	return svc.gitRepo.Storer.RemoveReference(ref)
}

//Add - adds the file content to the staging area
func (svc *service) Add(path string) error {

	wt, err := svc.gitRepo.Worktree()
	if err != nil {
		return nil
	}

	return wt.Add(path)
}

//Log - Gets the HEAD history from HEAD, just like command "git log"
func (svc *service) Log() ([]contract.Commit, error) {

	ref, err := svc.gitRepo.Head()
	if err != nil {
		return nil, err
	}

	cIter, err := svc.gitRepo.Log(&git.LogOptions{From: ref.Hash()})
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
