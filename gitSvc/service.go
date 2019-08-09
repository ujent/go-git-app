package gitsvc

import (
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ujent/go-git-app/contract"
	"github.com/ujent/go-git-mysql/mysqlfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
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
	//Clone() error
	Fetch() error
	Pull() error
	Push() error
	Commit() error
	Merge(branch string) error
	Branches() ([]string, error)
	Checkout(branch string) error
	DeleteBranch(branch string) error
	Add() error
	Log() error
}

type service struct {
	user     *contract.Credentials
	settings *contract.ServerSettings
	gitRepo  *git.Repository
}

//New - create an instance of gitSvc
func New(user *contract.Credentials, s *contract.ServerSettings) (Service, error) {

	if user.Name == "" {
		return nil, errors.New("userName cannot be empty")
	}

	if user.Email == "" {
		return nil, errors.New("userEmail cannot be empty")
	}

	return &service{user: user, settings: s}, nil
}

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

func (svc *service) OpenRepository(name string) error {
	if name == "" {
		return errors.New("Repository name cannot be empty")
	}

	db, err := sqlx.Connect("mysql", svc.settings.GitConnStr)

	if err != nil {
		return err
	}

	defer db.Close()

	tables := []string{}
	table := gitPrefix + name

	err = db.Get(&tables, "SELECT table_name FROM information_schema.tables WHERE table_type = 'base table' AND table_name = ?", table)
	if err != nil {
		return err
	}

	if len(tables) == 0 {
		return svc.CreateRepository(name)
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

// func (svc *service) Clone() error {
// 	return nil
// }

func (svc *service) Repositories() ([]string, error) {
	db, err := sqlx.Connect("mysql", svc.settings.GitConnStr)

	if err != nil {
		return nil, err
	}

	tables := []string{}
	db.Select(&tables, "SELECT table_name FROM information_schema.tables WHERE table_type = 'base table' ORDER BY table_name ASC")

	db.Close()

	repos := []string{}
	for _, t := range tables {
		if strings.HasPrefix(t, gitPrefix) {
			repos = append(repos, strings.TrimPrefix(t, gitPrefix))
		}
	}

	return repos, nil
}

func (svc *service) Branches() ([]string, error) {
	iter, err := svc.gitRepo.Branches()
	if err != nil {
		return nil, err
	}

	branches := []string{}

	iter.ForEach(func(br *plumbing.Reference) error {
		if !br.Name().IsRemote() {
			branches = append(branches, br.Name().Short())
		}
		return nil
	})

	return branches, nil
}

func (svc *service) Fetch() error {

	headRef, err := svc.gitRepo.Head()
	if err != nil {
		return err
	}

	current := headRef.Name()

	err = svc.gitRepo.Fetch(&git.FetchOptions{
		RemoteName: current.Short(),
	})

	if err != nil {
		return err
	}

	return nil
}

func (svc *service) Pull() error {
	return nil
}

func (svc *service) Push() error {
	return nil
}

func (svc *service) Commit() error {
	return nil
}

func (svc *service) Merge(branch string) error {
	return nil
}

func (svc *service) Checkout(branch string) error {
	return nil
}

func (svc *service) DeleteBranch(branch string) error {
	return nil
}

func (svc *service) Add() error {
	return nil
}

func (svc *service) Log() error {
	return nil
}
