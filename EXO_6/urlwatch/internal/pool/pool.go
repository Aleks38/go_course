package pool

import (
	"context"
	"sync"
	"time"

	"urlwatch/internal/domain"
)

func Run(ctx context.Context, c domain.Checker, urls []string, concurrency int, timeout, globalTimeout time.Duration) []domain.CheckResult {
	if concurrency < 1 {
		concurrency = 1
	}

	ctx, cancel := context.WithTimeout(ctx, globalTimeout)
	defer cancel()

	jobs := make(chan int)
	results := make(chan indexedResult)

	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go worker(ctx, c, urls, timeout, jobs, results, &wg)
	}

	go func() {
		for idx := range urls {
			jobs <- idx
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]domain.CheckResult, len(urls))
	for r := range results {
		out[r.index] = r.result
	}
	return out
}

func worker(ctx context.Context, c domain.Checker, urls []string, timeout time.Duration, jobs <-chan int, results chan<- indexedResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for idx := range jobs {
		results <- indexedResult{index: idx, result: check(ctx, c, urls[idx], timeout)}
	}
}

func check(ctx context.Context, c domain.Checker, url string, timeout time.Duration) domain.CheckResult {
	select {
	case <-ctx.Done():
		return domain.CheckResult{URL: url, OK: false, Error: ctx.Err().Error()}
	default:
	}

	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Check(cctx, url)
}
