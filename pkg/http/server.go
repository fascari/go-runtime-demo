package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	port   string
}

func NewServer(port string) *Server {
	return &Server{
		router: mux.NewRouter(),
		port:   port,
	}
}

func (s *Server) Router() *mux.Router {
	return s.router
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
