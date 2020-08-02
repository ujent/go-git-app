package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"bitbucket.org/vishjosh/bipp-go-git"
	"bitbucket.org/vishjosh/bipp-go-git/experimental-app/contract"
	"bitbucket.org/vishjosh/bipp-go-git/experimental-app/gitsvc"
)

type server struct {
	settings *contract.ServerSettings
	logger   *log.Logger
	gitSvc   gitsvc.Service
}

func newServer(settings *contract.ServerSettings, logger *log.Logger, db *sqlx.DB) (*server, error) {

	gitSvc, err := gitsvc.New(settings, db)
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

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		//MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/switch", s.switchUser)
		})

		r.Route("/repositories", func(r chi.Router) {
			r.Get("/{user}", s.repositories)
			r.Get("/current/{user}", s.currentRepo)
			r.Get("/open", s.openRepository)
			r.Post("/", s.createRepository)
			r.Post("/clone", s.clone)
			r.Delete("/", s.deleteRepository)
		})

		r.Route("/branches", func(r chi.Router) {
			r.Get("/", s.branches)
			r.Post("/checkout", s.checkoutBranch)
			r.Post("/", s.createBranch)
			r.Delete("/", s.deleteBranch)
		})

		r.Route("/log", func(r chi.Router) {
			r.Get("/", s.logs)
		})

		r.Route("/files", func(r chi.Router) {
			r.Get("/all", s.files)
			r.Get("/", s.file)
			r.Post("/", s.addFile)
			r.Put("/", s.editFile)
			r.Delete("/", s.removeFile)
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
			r.Post("/abort", s.abortMerge)
		})
	})

	err := http.ListenAndServe(":"+s.settings.Port, r)
	if err != nil {
		return err
	}

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

	err = s.gitSvc.SwitchUser(&contract.User{Name: rq.Name})

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
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
		msg, err = s.gitSvc.Pull(s.toBaseRequest(rq.Base), rq.Remote, nil)

	} else {
		msg, err = s.gitSvc.Pull(s.toBaseRequest(rq.Base), rq.Remote, &contract.Credentials{Name: rq.Auth.Name, Password: rq.Auth.Psw})
	}

	if err != nil {
		if msg != "" {
			s.writeJSON(w, http.StatusOK, &contract.MsgResult{Msg: msg})
		} else {
			s.writeError(w, http.StatusInternalServerError, err)
		}

		return
	}

	s.writeJSON(w, http.StatusOK, &contract.MsgResult{Msg: msg})
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
		err = s.gitSvc.Push(s.toBaseRequest(rq.Base), rq.Remote, nil)

	} else {
		err = s.gitSvc.Push(s.toBaseRequest(rq.Base), rq.Remote, &contract.Credentials{Name: rq.Auth.Name, Password: rq.Auth.Psw})
	}

	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) commit(w http.ResponseWriter, r *http.Request) {
	rq := &contract.CommitRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	_, err = s.gitSvc.Commit(s.toBaseRequest(rq.Base), rq.Message)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) toBaseRequest(rq *contract.BaseRequestRQ) *contract.BaseRequest {
	if rq == nil {
		return nil
	}

	return &contract.BaseRequest{Repository: rq.Repository, Branch: rq.Branch, User: &contract.User{Name: rq.User}}
}

func (s *server) clone(w http.ResponseWriter, r *http.Request) {
	rq := &contract.CloneRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	var repo string

	if rq.Auth == nil {
		repo, err = s.gitSvc.Clone(rq.User, rq.URL, nil)

	} else {
		repo, err = s.gitSvc.Clone(rq.User, rq.URL, &contract.Credentials{Name: rq.Auth.Name, Password: rq.Auth.Psw})
	}

	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	s.writeJSON(w, http.StatusOK, &contract.RepoRS{Name: repo})
}

func (s *server) files(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	branch := q.Get("branch")

	if branch == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("branch cannot be empty"))

		return
	}

	repo := q.Get("repo")

	if repo == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("repo cannot be empty"))

		return
	}

	user := q.Get("user")

	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user cannot be empty"))

		return
	}

	files, err := s.gitSvc.FilesList(&contract.BaseRequest{User: &contract.User{Name: user}, Repository: repo, Branch: branch})
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	statuses, err := s.gitSvc.Status(&contract.BaseRequest{User: &contract.User{Name: user}, Repository: repo, Branch: branch})
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	res := []contract.FileInfoRS{}
	for _, f := range files {

		st, ok := statuses[f.Path]

		if ok {
			res = append(res, contract.FileInfoRS{Path: f.Path, IsConflict: f.IsConflict, FileStatus: s.toFileStatus(st.Staging)})
			delete(statuses, f.Path)
		} else {
			res = append(res, contract.FileInfoRS{Path: f.Path, IsConflict: f.IsConflict, FileStatus: contract.UnmodifiedFileStatus})
		}
	}

	for path, st := range statuses {
		res = append(res, contract.FileInfoRS{Path: path, FileStatus: s.toFileStatus(st.Staging)})
	}

	s.writeJSON(w, http.StatusOK, &contract.FilesRS{Files: res})
}

func (s *server) toFileStatus(c git.StatusCode) contract.FileStatus {
	switch c {
	case git.Unmodified:
		return contract.UnmodifiedFileStatus
	case git.Untracked:
		return contract.UntrackedFileStatus
	case git.Modified:
		return contract.ModifiedFileStatus
	case git.Added:
		return contract.AddedFileStatus
	case git.Deleted:
		return contract.DeletedFileStatus
	case git.Renamed:
		return contract.RenamedFileStatus
	case git.Copied:
		return contract.CopiedFileStatus
	case git.UpdatedButUnmerged:
		return contract.UpdatedButUnmergedFileStatus
	default:
		return contract.UnspecifiedFileStatus
	}
}

func (s *server) file(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	branch := q.Get("branch")

	if branch == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("branch cannot be empty"))

		return
	}

	repo := q.Get("repo")

	if repo == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("repo cannot be empty"))

		return
	}

	user := q.Get("user")

	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user cannot be empty"))

		return
	}

	path := q.Get("path")

	if path == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("path cannot be empty"))

		return
	}

	f, err := s.gitSvc.File(&contract.BaseRequest{User: &contract.User{Name: user}, Repository: repo, Branch: branch}, path)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil && err != io.EOF {
		s.writeError(w, http.StatusInternalServerError, err)
	}

	s.writeJSON(w, http.StatusOK, &contract.FileRS{Path: path, Content: string(bytes)})
}

func (s *server) addFile(w http.ResponseWriter, r *http.Request) {
	rq := &contract.AddFileRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	err = s.gitSvc.AddFile(s.toBaseRequest(rq.Base), rq.Path, rq.Content)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) editFile(w http.ResponseWriter, r *http.Request) {
	rq := &contract.EditFileRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	err = s.gitSvc.EditFile(s.toBaseRequest(rq.Base), rq.Path, rq.Content)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) removeFile(w http.ResponseWriter, r *http.Request) {
	rq := &contract.RemoveFileRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	err = s.gitSvc.RemoveFile(s.toBaseRequest(rq.Base), rq.Path)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) logs(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	branch := q.Get("branch")

	if branch == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("branch cannot be empty"))

		return
	}

	repo := q.Get("repo")

	if repo == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("repo cannot be empty"))

		return
	}

	user := q.Get("user")

	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user cannot be empty"))

		return
	}

	commits, err := s.gitSvc.Log(&contract.BaseRequest{User: &contract.User{Name: user}, Repository: repo, Branch: branch})
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

	err = s.gitSvc.CreateBranch(rq.User, rq.Repo, rq.Branch, "")
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) deleteBranch(w http.ResponseWriter, r *http.Request) {
	rq := &contract.BranchRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	err = s.gitSvc.RemoveBranch(rq.User, rq.Repo, rq.Branch)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) checkoutBranch(w http.ResponseWriter, r *http.Request) {
	rq := &contract.BranchRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	err = s.gitSvc.CheckoutBranch(rq.User, rq.Repo, rq.Branch)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) branches(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	repo := q.Get("repo")

	if repo == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("repo cannot be empty"))

		return
	}

	user := q.Get("user")

	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user cannot be empty"))

		return
	}

	branches, err := s.gitSvc.Branches(user, repo)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	cur := ""
	if len(branches) > 0 {
		cBr, err := s.gitSvc.CurrentBranch()
		if err != nil {
			s.writeError(w, http.StatusInternalServerError, err)

			return
		}

		if cBr != nil {
			cur = cBr.Name
		}
	}

	s.writeJSON(w, http.StatusOK, &contract.BranchesRS{Branches: branches, Current: cur})
}

func (s *server) createRepository(w http.ResponseWriter, r *http.Request) {

	rq := &contract.RepoRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if rq.Repo == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("repo cannot be empty"))

		return
	}

	if rq.User == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user cannot be empty"))

		return
	}

	err = s.gitSvc.CreateRepository(rq.User, rq.Repo)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) deleteRepository(w http.ResponseWriter, r *http.Request) {
	rq := &contract.RepoRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if rq.Repo == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("repo cannot be empty"))

		return
	}

	if rq.User == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user cannot be empty"))

		return
	}

	err = s.gitSvc.RemoveRepository(rq.User, rq.Repo)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) openRepository(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	repo := q.Get("repo")

	if repo == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("repo cannot be empty"))

		return
	}

	user := q.Get("user")

	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user cannot be empty"))

		return
	}

	err := s.gitSvc.OpenRepository(user, repo)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *server) currentRepo(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")

	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user cannot be empty"))

		return
	}

	repo := s.gitSvc.CurrentRepository()

	s.writeJSON(w, http.StatusOK, &contract.RepoRS{Name: repo})
}

func (s *server) repositories(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")

	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user cannot be empty"))

		return
	}

	repos, err := s.gitSvc.Repositories(user)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)

		return
	}

	res := []string{}

	for _, r := range repos {

		res = append(res, r)
	}

	s.writeJSON(w, http.StatusOK, &contract.RepositoriesRS{Repos: res})
}

func (s *server) merge(w http.ResponseWriter, r *http.Request) {
	rq := &contract.MergeRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if rq.Theirs == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Theirs branch cannot be empty"))

		return
	}

	mergeMsg, err := s.gitSvc.Merge(s.toBaseRequest(rq.Base), rq.Theirs)

	if err != nil {

		switch err {
		case git.ErrHasUncommittedFiles:
			{
				//ErrHasUncommittedFiles occurs when there are any unstaged or staged files before merge
				s.writeJSON(w, http.StatusOK, &contract.MergeRS{Message: err.Error()})
				return
			}
		case git.ErrMergeCommitNeeded:
			{
				s.writeJSON(w, http.StatusOK, &contract.MergeRS{Message: err.Error()})
				return
			}
		case git.ErrMergeWithConflicts:
			{
				s.writeJSON(w, http.StatusOK, &contract.MergeRS{Message: mergeMsg})
				return

			}
		default:
			{
				s.writeError(w, http.StatusInternalServerError, err)
				return
			}
		}
	}

	s.writeJSON(w, http.StatusOK, &contract.MergeRS{IsFastforward: true})
}

func (s *server) abortMerge(w http.ResponseWriter, r *http.Request) {
	rq := &contract.AbortMergeRQ{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(rq)

	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	err = s.gitSvc.AbortMerge(s.toBaseRequest(rq.Base))

	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
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
