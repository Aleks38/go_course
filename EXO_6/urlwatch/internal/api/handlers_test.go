package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"urlwatch/internal/domain"
)

type mockChecker struct{}

func (mockChecker) Check(_ context.Context, url string) domain.CheckResult {
	return domain.CheckResult{URL: url, StatusCode: 200, OK: true}
}

type mockStore struct {
	saved map[string]domain.Batch
}

func newMockStore() *mockStore {
	return &mockStore{saved: make(map[string]domain.Batch)}
}

func (m *mockStore) Save(_ context.Context, b domain.Batch) error {
	m.saved[b.ID] = b
	return nil
}

func (m *mockStore) Get(_ context.Context, id string) (domain.Batch, error) {
	b, ok := m.saved[id]
	if !ok {
		return domain.Batch{}, domain.ErrBatchNotFound
	}
	return b, nil
}

func TestHandleChecksSuccess(t *testing.T) {
	srv := New(mockChecker{}, newMockStore(), nil)

	body := `{"urls":["https://go.dev","https://example.test"]}`
	req := httptest.NewRequest(http.MethodPost, "/v1/checks", strings.NewReader(body))
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201 ; body = %s", rec.Code, rec.Body.String())
	}

	var got domain.Batch
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("reponse JSON invalide: %v", err)
	}
	if got.ID == "" {
		t.Error("batch_id manquant")
	}
	if got.Summary.Total != 2 || got.Summary.Up != 2 {
		t.Errorf("summary inattendu: %+v", got.Summary)
	}
	if len(got.Results) != 2 {
		t.Errorf("results = %d, want 2", len(got.Results))
	}
}

func TestHandleGetCheckNotFound(t *testing.T) {
	srv := New(mockChecker{}, newMockStore(), nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/checks/b_inconnu", nil)
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 ; body = %s", rec.Code, rec.Body.String())
	}

	var got errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("reponse JSON invalide: %v", err)
	}
	if got.Error.Code != "batch_not_found" {
		t.Errorf("code = %q, want batch_not_found", got.Error.Code)
	}
}

func TestHandleChecksValidationError(t *testing.T) {
	srv := New(mockChecker{}, newMockStore(), nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/checks", strings.NewReader(`{"urls":[]}`))
	rec := httptest.NewRecorder()

	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 ; body = %s", rec.Code, rec.Body.String())
	}

	var got errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("reponse JSON invalide: %v", err)
	}
	if got.Error.Code != "invalid_request" {
		t.Errorf("code = %q, want invalid_request", got.Error.Code)
	}
}
