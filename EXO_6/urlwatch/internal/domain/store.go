package domain

import "context"

type Store interface {
	Save(ctx context.Context, b Batch) error
	Get(ctx context.Context, id string) (Batch, error)
}
