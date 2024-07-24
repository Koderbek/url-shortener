package server

import (
	"github.com/Koderbek/url-shortener/internal/app/config"
	"github.com/Koderbek/url-shortener/internal/app/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type server struct {
	router *chi.Mux
}

func newServer() *server {
	s := &server{router: chi.NewRouter()}
	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/api/shorten", logger.RequestLogger(apiShorten))
	s.router.HandleFunc("/", logger.RequestLogger(shorten))
	s.router.HandleFunc("/{id}", logger.RequestLogger(findURL))
}

func Start() {
	if err := logger.Initialize(config.Config.Flags.LogLevel); err != nil {
		panic(err)
	}

	s := newServer()
	err := http.ListenAndServe(config.Config.Flags.ServerAddress, s.router)
	if err != nil {
		panic(err)
	}
}
