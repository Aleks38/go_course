package api

import (
	"fmt"
	"strings"
	"time"

	"urlwatch/internal/domain"
)

const (
	maxURLs = 100

	defaultConcurrency = 8
	minConcurrency     = 1
	maxConcurrency     = 50

	defaultTimeoutMs = 5000
	minTimeoutMs     = 100
	maxTimeoutMs     = 30000
)

func validate(req createChecksRequest) (urls []string, concurrency int, timeout time.Duration, err error) {
	if len(req.URLs) == 0 {
		return nil, 0, 0, &domain.ValidationError{Field: "urls", Message: "au moins une URL est requise"}
	}
	if len(req.URLs) > maxURLs {
		return nil, 0, 0, &domain.ValidationError{Field: "urls", Message: fmt.Sprintf("%d URLs maximum", maxURLs)}
	}
	for _, u := range req.URLs {
		if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
			return nil, 0, 0, &domain.ValidationError{Field: "urls", Message: fmt.Sprintf("URL invalide : %q", u)}
		}
	}

	concurrency = defaultConcurrency
	timeoutMs := defaultTimeoutMs
	if req.Options != nil {
		if req.Options.Concurrency != nil {
			concurrency = *req.Options.Concurrency
			if concurrency < minConcurrency || concurrency > maxConcurrency {
				return nil, 0, 0, &domain.ValidationError{Field: "options.concurrency", Message: fmt.Sprintf("doit être compris entre %d et %d", minConcurrency, maxConcurrency)}
			}
		}
		if req.Options.TimeoutMs != nil {
			timeoutMs = *req.Options.TimeoutMs
			if timeoutMs < minTimeoutMs || timeoutMs > maxTimeoutMs {
				return nil, 0, 0, &domain.ValidationError{Field: "options.timeout_ms", Message: fmt.Sprintf("doit être compris entre %d et %d", minTimeoutMs, maxTimeoutMs)}
			}
		}
	}

	return req.URLs, concurrency, time.Duration(timeoutMs) * time.Millisecond, nil
}
