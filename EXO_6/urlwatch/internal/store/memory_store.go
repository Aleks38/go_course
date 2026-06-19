package store

import (
	"context"
	"fmt"
	"sync"

	"urlwatch/internal/domain"
)

type MemoryStore struct {
	mu      sync.Mutex
	batches map[string]domain.Batch
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{batches: make(map[string]domain.Batch)}
}

func (s *MemoryStore) Save(ctx context.Context, b domain.Batch) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.batches[b.ID] = b
	return nil
}

func (s *MemoryStore) Get(ctx context.Context, id string) (domain.Batch, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.batches[id]
	if !ok {
		return domain.Batch{}, fmt.Errorf("get %s: %w", id, domain.ErrBatchNotFound)
	}
	return b, nil
}
