package api

import (
	"errors"
	"testing"
	"time"

	"urlwatch/internal/domain"
)

func intPtr(i int) *int { return &i }

func TestValidate(t *testing.T) {
	manyURLs := make([]string, 101)
	for i := range manyURLs {
		manyURLs[i] = "https://example"
	}

	tests := []struct {
		name            string
		req             createChecksRequest
		wantErr         bool
		wantField       string
		wantConcurrency int
		wantTimeout     time.Duration
	}{
		{
			name:            "valide avec defauts",
			req:             createChecksRequest{URLs: []string{"https://go.dev"}},
			wantConcurrency: 8,
			wantTimeout:     5000 * time.Millisecond,
		},
		{
			name: "valide avec options",
			req: createChecksRequest{
				URLs:    []string{"http://example.com"},
				Options: &createChecksOptions{Concurrency: intPtr(4), TimeoutMs: intPtr(2000)},
			},
			wantConcurrency: 4,
			wantTimeout:     2000 * time.Millisecond,
		},
		{
			name:      "urls vide",
			req:       createChecksRequest{URLs: nil},
			wantErr:   true,
			wantField: "urls",
		},
		{
			name:      "trop d urls",
			req:       createChecksRequest{URLs: manyURLs},
			wantErr:   true,
			wantField: "urls",
		},
		{
			name:      "url non http",
			req:       createChecksRequest{URLs: []string{"ftp://example.com"}},
			wantErr:   true,
			wantField: "urls",
		},
		{
			name:      "concurrency hors borne",
			req:       createChecksRequest{URLs: []string{"https://go.dev"}, Options: &createChecksOptions{Concurrency: intPtr(0)}},
			wantErr:   true,
			wantField: "options.concurrency",
		},
		{
			name:      "timeout hors borne",
			req:       createChecksRequest{URLs: []string{"https://go.dev"}, Options: &createChecksOptions{TimeoutMs: intPtr(50)}},
			wantErr:   true,
			wantField: "options.timeout_ms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls, concurrency, timeout, err := validate(tt.req)

			if tt.wantErr {
				var ve *domain.ValidationError
				if !errors.As(err, &ve) {
					t.Fatalf("erreur attendue (*domain.ValidationError), got %v", err)
				}
				if ve.Field != tt.wantField {
					t.Errorf("Field = %q, want %q", ve.Field, tt.wantField)
				}
				return
			}

			if err != nil {
				t.Fatalf("erreur inattendue: %v", err)
			}
			if len(urls) != len(tt.req.URLs) {
				t.Errorf("urls len = %d, want %d", len(urls), len(tt.req.URLs))
			}
			if concurrency != tt.wantConcurrency {
				t.Errorf("concurrency = %d, want %d", concurrency, tt.wantConcurrency)
			}
			if timeout != tt.wantTimeout {
				t.Errorf("timeout = %v, want %v", timeout, tt.wantTimeout)
			}
		})
	}
}
