package main

import (
	"encoding/json"
	"fmt"
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

	r.Route("/users", func(r chi.Router) {
		r.Post("/switch", s.switchUser)
	})

	r.Route("/repositories", func(r chi.Router) {
		r.Get("/", s.repositories)
		r.Get("/open/{name}", s.openRepository)
		r.Post("/", s.createRepository)
		r.Post("/clone", s.clone)
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
		r.Get("/", s.files)
	})

	r.Route("/commit", func(r chi.Router) {
		r.Post("/", s.commit)
	})

	r.Route("/pull", func(r chi.Router) {
		r.Post("/", s.pull)
	})

	r.Route("/push", func(r chi.Router) {
		r.Post("/", s.push)
	})

	r.Route("/merge", func(r chi.Router) {
		r.Post("/", s.merge)
	})

	return nil
}

func (s *server) switchUser(w http.ResponseWriter, r *http.Request) {
	rq := &contract.SwitchUserRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if rq.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name cannot be empty"))

		return
	}

	var user *contract.User

	for _, u := range contract.TestUsers {
		if u.Name == rq.Name {
			user = &u
			break
		}
	}

	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("User %s doesn't exist", rq.Name)))

		return
	}

	err = s.gitSvc.SwitchUser(user)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) pull(w http.ResponseWriter, r *http.Request) {
	rq := &contract.PullRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	var msg string
	if rq.Auth == nil {
		msg, err = s.gitSvc.Pull(rq.Remote, nil)

	} else {
		msg, err = s.gitSvc.Pull(rq.Remote, &contract.Credentials{Name: rq.Auth.Name, Password: rq.Auth.Psw})
	}

	if err != nil {
		if msg != "" {
			s.writeJSON(w, http.StatusOK, &contract.MsgResult{Msg: msg})
		} else {
			s.writeError(w, http.StatusInternalServerError, err)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) push(w http.ResponseWriter, r *http.Request) {
	rq := &contract.PushRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if rq.Auth == nil {
		err = s.gitSvc.Push(rq.Remote, nil)

	} else {
		err = s.gitSvc.Push(rq.Remote, &contract.Credentials{Name: rq.Auth.Name, Password: rq.Auth.Psw})
	}

	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) commit(w http.ResponseWriter, r *http.Request) {
	rq := &contract.CommitRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	_, err = s.gitSvc.Commit(rq.Message)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) clone(w http.ResponseWriter, r *http.Request) {
	rq := &contract.CloneRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if rq.Auth == nil {
		err = s.gitSvc.Clone(rq.URL, rq.RepoName, nil)

	} else {
		err = s.gitSvc.Clone(rq.URL, rq.RepoName, &contract.Credentials{Name: rq.Auth.Name, Password: rq.Auth.Psw})
	}

	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) files(w http.ResponseWriter, r *http.Request) {
	files, err := s.gitSvc.FilesList()
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
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
		return
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
		return
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
		return
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
	var cur string

	for _, b := range branches {
		if b == n.Name {
			res = append(res, contract.BranchRS{Name: b})
			cur = b
		} else {
			res = append(res, contract.BranchRS{Name: n.Name})
		}
	}

	s.writeJSON(w, http.StatusOK, &contract.BranchesRS{Branches: res, Current: cur})
}

func (s *server) createRepository(w http.ResponseWriter, r *http.Request) {

	repo := &contract.RepoRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(repo)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
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
		return
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
	var cur string

	for _, r := range repos {
		if r == n {
			res = append(res, contract.RepoRS{Name: r})
			cur = r
		} else {
			res = append(res, contract.RepoRS{Name: r})
		}
	}

	s.writeJSON(w, http.StatusOK, &contract.RepositoriesRS{Repos: res, Current: cur})
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

func (s *server) merge(w http.ResponseWriter, r *http.Request) {

}

/*
func (s *server) merge(w http.ResponseWriter, r *http.Request) {
	rq := &contract.MergeRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if rq.Branch == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("branch cannot be empty"))

		return
	}

	err = s.gitSvc.Merge(rq.Branch)

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
*/
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
