package api

import (
	"net/http"

	"urlwatch/internal/domain"
)

type Server struct {
	checker domain.Checker
	store   domain.Store
}

func New(checker domain.Checker, store domain.Store) *Server {
	return &Server{checker: checker, store: store}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/checks", s.handleChecks)
	mux.HandleFunc("/v1/checks/", s.handleCheckByID)
	mux.HandleFunc("/healthz", s.handleHealthz)
	return mux
}
