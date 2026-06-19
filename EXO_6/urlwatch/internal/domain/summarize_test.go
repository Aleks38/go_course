package domain

import (
	"testing"
	"time"
)

func TestSummarize(t *testing.T) {
	tests := []struct {
		name     string
		results  []CheckResult
		duration time.Duration
		want     Summary
	}{
		{
			name:     "lot vide",
			results:  nil,
			duration: 0,
			want:     Summary{Total: 0, Up: 0, Down: 0, DurationMs: 0},
		},
		{
			name:     "melange up et down",
			results:  []CheckResult{{OK: true}, {OK: false}, {OK: true}},
			duration: 1500 * time.Millisecond,
			want:     Summary{Total: 3, Up: 2, Down: 1, DurationMs: 1500},
		},
		{
			name:     "tous up",
			results:  []CheckResult{{OK: true}, {OK: true}},
			duration: 200 * time.Millisecond,
			want:     Summary{Total: 2, Up: 2, Down: 0, DurationMs: 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Summarize(tt.results, tt.duration)
			if got != tt.want {
				t.Errorf("Summarize() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
