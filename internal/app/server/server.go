package server

import (
	"github.com/Koderbek/url-shortener/internal/app/config"
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
	s.router.HandleFunc("/", shorten)
	s.router.HandleFunc("/{id}", findURL)
}

func Start() {
	s := newServer()
	err := http.ListenAndServe(config.Config.Flags.ServerAddress, s.router)
	if err != nil {
		panic(err)
	}
}
