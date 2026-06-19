package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"urlwatch/internal/domain"
	"urlwatch/internal/pool"
)

const globalTimeout = 60 * time.Second

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "méthode non autorisée")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleChecks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "méthode non autorisée")
		return
	}

	var req createChecksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "corps JSON invalide")
		return
	}

	urls, concurrency, timeout, err := validate(req)
	if err != nil {
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			writeError(w, http.StatusBadRequest, "invalid_request", ve.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal", "erreur interne")
		return
	}

	start := time.Now()
	results := pool.Run(r.Context(), s.checker, urls, concurrency, timeout, globalTimeout)
	batch := domain.Batch{
		ID:        newBatchID(),
		CreatedAt: time.Now().UTC(),
		Summary:   domain.Summarize(results, time.Since(start)),
		Results:   results,
	}
	setBatchID(r.Context(), batch.ID)

	if err := s.store.Save(r.Context(), batch); err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "echec de la persistance du lot")
		return
	}

	writeJSON(w, http.StatusCreated, batch)
}

func (s *Server) handleCheckByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "méthode non autorisée")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/checks/")
	setBatchID(r.Context(), id)

	batch, err := s.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrBatchNotFound) {
			writeError(w, http.StatusNotFound, "batch_not_found", fmt.Sprintf("aucun lot avec l'id %s", id))
			return
		}
		writeError(w, http.StatusInternalServerError, "internal", "erreur interne")
		return
	}

	writeJSON(w, http.StatusOK, batch)
}
