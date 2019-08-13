package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ujent/go-git-app/contract"
	"github.com/ujent/go-git-app/gitsvc"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
)

type server struct {
	settings *contract.ServerSettings
	logger   *log.Logger
	gitSvc   gitsvc.Service
}

func newServer(settings *contract.ServerSettings, user *contract.User, logger *log.Logger) (*server, error) {
	db, err := sqlx.Connect("mysql", settings.GitConnStr)

	if err != nil {
		return nil, err
	}

	defer db.Close()

	gitSvc, err := gitsvc.New(user, settings, db)
	if err != nil {
		return nil, err
	}

	s := server{logger: logger, settings: settings, gitSvc: gitSvc}

	return &s, nil
}

func (s *server) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	err := http.ListenAndServe(":"+s.settings.Port, r)
	if err != nil {
		return err
	}
	//r.HandleFunc()
	//router := http.NewServeMux()
	//router.HandleFunc(, s.)

	return nil
}

func (s *server) handleMergeCommit(msg string) error {
	_, err := s.gitSvc.Commit(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) handleConflictFiles(path string) ([]contract.MergeFile, error) {
	res, err := s.gitSvc.ConflictFiles(path)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) handleConflictResultFile(path string) (billy.File, error) {
	f, err := s.gitSvc.ConflictResultFile(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (s *server) handleMerge(branch string) (string, error) {
	err := s.gitSvc.Merge(branch)

	if err != nil {

		switch err {
		case git.ErrHasUncommittedFiles:
			{
				//ErrHasUncommittedFiles occurs when there are any unstaged or staged files before merge
				return "", err
			}
		case git.ErrMergeCommitNeeded:
			{
				msg, err := s.gitSvc.MergeMsgShort()
				if err != nil {
					return "", err
				}

				return msg, nil
			}
		case git.ErrMergeWithConflicts:
			{
				//ToDo RA - maybe write to error instead of Stdout?

				msg, err := s.gitSvc.MergeMsgFull()
				if err != nil {
					return "", err
				}

				//ToDo RA - write to Response info about conflicts

				//conflictFiles := s.svc.ConflictFileList()
				return msg, nil

			}
		default:
			{
				return "", err
			}
		}

	} else {
		//if there is no error it means that it was a fastforward merge

		//ToDo RA - write that all is OK or that it was successful ff merge
		return "", nil
	}
}
