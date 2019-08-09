package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ujent/go-git-app/contract"
	gitsvc "github.com/ujent/go-git-app/gitSvc"
)

type server struct {
	settings *contract.ServerSettings
	logger   *log.Logger
	gitSvc   gitsvc.Service
}

func newServer(settings *contract.ServerSettings, user *contract.Credentials, logger *log.Logger) (*server, error) {
	gitSvc, err := gitsvc.New(user, settings)
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
