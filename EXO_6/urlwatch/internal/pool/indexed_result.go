package pool

import "urlwatch/internal/domain"

type indexedResult struct {
	index  int
	result domain.CheckResult
}
