package domain

import "context"

type Checker interface {
	Check(ctx context.Context, url string) CheckResult
}
