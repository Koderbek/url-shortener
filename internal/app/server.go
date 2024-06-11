package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

type server struct {
	router *mux.Router
}

func newServer() *server {
	s := &server{router: mux.NewRouter()}
	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/", shorten)
	s.router.HandleFunc("/{id}", findURL)
}

func Start() {
	s := newServer()
	err := http.ListenAndServe(`:8080`, s.router)
	if err != nil {
		panic(err)
	}
}
