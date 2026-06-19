package domain

import "time"

type Batch struct {
	ID        string        `json:"batch_id"`
	CreatedAt time.Time     `json:"created_at"`
	Summary   Summary       `json:"summary"`
	Results   []CheckResult `json:"results"`
}
