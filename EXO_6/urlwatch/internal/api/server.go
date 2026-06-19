package api

import (
	"log/slog"
	"net/http"

	"urlwatch/internal/domain"
)

type Server struct {
	checker domain.Checker
	store   domain.Store
	logger  *slog.Logger
}

func New(checker domain.Checker, store domain.Store, logger *slog.Logger) *Server {
	return &Server{checker: checker, store: store, logger: logger}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/checks", s.handleChecks)
	mux.HandleFunc("/v1/checks/", s.handleCheckByID)
	mux.HandleFunc("/healthz", s.handleHealthz)
	if s.logger == nil {
		return mux
	}
	return loggingMiddleware(s.logger, mux)
}
