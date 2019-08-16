package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ujent/go-git-app/contract"
	"github.com/ujent/go-git-app/gitsvc"
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

	r.Route("/repositories", func(r chi.Router) {
		r.Get("/", s.repositories)
		r.Get("/open/{name}", s.openRepository)
		r.Post("/", s.createRepository)
		r.Delete("/", s.deleteRepository)
	})

	r.Route("/branches", func(r chi.Router) {
		r.Get("/", s.branches)
		r.Get("/checkout/{name}", s.checkoutBranch)
		r.Post("/", s.createBranch)
		r.Delete("/", s.deleteBranch)
	})

	r.Route("/log", func(r chi.Router) {
		r.Get("/", s.logs)
	})

	r.Route("/files", func(r chi.Router) {
		r.Get("/all", s.files)
	})

	return nil
}

func (s *server) files(w http.ResponseWriter, r *http.Request) {
	files, err := s.gitSvc.FilesList()
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
	}

	res := []contract.FileInfoRS{}
	for _, f := range files {
		res = append(res, contract.FileInfoRS{Path: f.Path, IsConflict: f.IsConflict})
	}

	s.writeJSON(w, http.StatusOK, &contract.FilesRS{Files: res})
}

func (s *server) logs(w http.ResponseWriter, r *http.Request) {
	commits, err := s.gitSvc.Log()
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
	}

	res := []contract.CommitRS{}

	for _, c := range commits {
		res = append(res, s.toCommitRS(c))
	}

	s.writeJSON(w, http.StatusOK, &contract.LogRS{Commits: res})
}

func (s *server) toCommitRS(c contract.Commit) contract.CommitRS {

	res := contract.CommitRS{
		Hash:    c.Hash,
		Message: c.Message,
		Date:    c.Date,
	}

	if c.Author != nil {
		res.Author = &contract.UserRS{Name: c.Author.Name, Email: c.Author.Email}
	}

	return res
}

func (s *server) createBranch(w http.ResponseWriter, r *http.Request) {
	rq := &contract.BranchRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
	}

	if rq.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name cannot be empty"))

		return
	}

	err = s.gitSvc.CreateBranch(rq.Name, "")
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) deleteBranch(w http.ResponseWriter, r *http.Request) {
	rq := &contract.BranchRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
	}

	if rq.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name cannot be empty"))

		return
	}

	err = s.gitSvc.RemoveBranch(rq.Name)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) checkoutBranch(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name cannot be empty"))

		return
	}

	err := s.gitSvc.CheckoutBranch(name)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) branches(w http.ResponseWriter, r *http.Request) {
	branches, err := s.gitSvc.Branches()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	n, err := s.gitSvc.CurrentBranch()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	res := []contract.BranchRS{}

	for _, b := range branches {
		if b == n.Name {
			res = append(res, contract.BranchRS{Name: b, IsCurrent: true})
		} else {
			res = append(res, contract.BranchRS{Name: n.Name, IsCurrent: false})
		}
	}

	s.writeJSON(w, http.StatusOK, &contract.BranchesRS{Branches: res})
}

func (s *server) createRepository(w http.ResponseWriter, r *http.Request) {

	repo := &contract.RepoRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(repo)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
	}

	if repo.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name cannot be empty"))

		return
	}

	err = s.gitSvc.CreateRepository(repo.Name)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) deleteRepository(w http.ResponseWriter, r *http.Request) {
	repo := &contract.RepoRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(repo)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
	}

	if repo.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name cannot be empty"))

		return
	}

	err = s.gitSvc.RemoveRepository(repo.Name)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) openRepository(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name cannot be empty"))

		return
	}

	err := s.gitSvc.OpenRepository(name)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) repositories(w http.ResponseWriter, r *http.Request) {
	repos, err := s.gitSvc.Repositories()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	n := s.gitSvc.CurrentRepository()
	res := []contract.RepoRS{}

	for _, r := range repos {
		if r == n {
			res = append(res, contract.RepoRS{Name: r, IsCurrent: true})
		} else {
			res = append(res, contract.RepoRS{Name: r, IsCurrent: false})
		}
	}

	s.writeJSON(w, http.StatusOK, &contract.RepositoriesRS{Repos: res})
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

func (s *server) handleConflictResultFile(path string) (string, error) {
	f, err := s.gitSvc.ConflictResultFile(path)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil && err != io.EOF {
		return "", err
	}

	return string(bytes), nil
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

func (s *server) writeJSON(w http.ResponseWriter, statusCode int, payload interface{}) {

	json, err := json.Marshal(payload)
	if err != nil {
		s.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(json)
}

func (s *server) writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}
