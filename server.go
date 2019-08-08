package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ujent/go-git-app/contract"
)

type server struct {
	settings *contract.ServerSettings
	logger   *log.Logger
}

func newServer(settings *contract.ServerSettings, logger *log.Logger) *server {
	s := server{logger: logger, settings: settings}

	return &s
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
