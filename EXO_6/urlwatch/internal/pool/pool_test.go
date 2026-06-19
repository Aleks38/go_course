package pool

import (
	"context"
	"sync"
	"testing"
	"time"

	"urlwatch/internal/domain"
)

type mockChecker struct {
	delay   time.Duration
	mu      sync.Mutex
	current int
	maxSeen int
}

func (m *mockChecker) Check(ctx context.Context, url string) domain.CheckResult {
	m.mu.Lock()
	m.current++
	if m.current > m.maxSeen {
		m.maxSeen = m.current
	}
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		m.current--
		m.mu.Unlock()
	}()

	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return domain.CheckResult{URL: url, OK: false, Error: ctx.Err().Error()}
		}
	}
	return domain.CheckResult{URL: url, StatusCode: 200, OK: true}
}

func (m *mockChecker) maxConcurrent() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.maxSeen
}

func TestRunReturnsResultForEachURLInOrder(t *testing.T) {
	urls := []string{"https://a", "https://b", "https://c"}
	got := Run(context.Background(), &mockChecker{}, urls, 2, time.Second, time.Second)

	if len(got) != len(urls) {
		t.Fatalf("len = %d, want %d", len(got), len(urls))
	}
	for i, r := range got {
		if r.URL != urls[i] {
			t.Errorf("ordre non preserve: got[%d].URL = %q, want %q", i, r.URL, urls[i])
		}
		if !r.OK {
			t.Errorf("resultat %d devrait etre OK", i)
		}
	}
}

func TestRunRespectsConcurrencyLimit(t *testing.T) {
	urls := make([]string, 20)
	for i := range urls {
		urls[i] = "https://example"
	}
	m := &mockChecker{delay: 20 * time.Millisecond}
	concurrency := 4

	Run(context.Background(), m, urls, concurrency, time.Second, time.Second)

	if got := m.maxConcurrent(); got > concurrency {
		t.Errorf("concurrence max observee = %d, ne doit pas depasser %d", got, concurrency)
	}
}

func TestRunHonorsCancellation(t *testing.T) {
	urls := []string{"https://a", "https://b"}
	m := &mockChecker{delay: time.Hour}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	got := Run(ctx, m, urls, 2, time.Second, time.Second)

	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("annulation non respectee (a pris %v)", elapsed)
	}
	for _, r := range got {
		if r.OK {
			t.Errorf("resultat devrait etre en echec apres annulation: %+v", r)
		}
	}
}

func TestRunHonorsPerURLTimeout(t *testing.T) {
	urls := []string{"https://slow"}
	m := &mockChecker{delay: time.Hour}

	start := time.Now()
	got := Run(context.Background(), m, urls, 1, 50*time.Millisecond, time.Second)

	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("timeout par URL non respecte (a pris %v)", elapsed)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].OK {
		t.Errorf("resultat devrait etre en echec apres timeout: %+v", got[0])
	}
}
