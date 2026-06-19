package domain

import "time"

func Summarize(results []CheckResult, duration time.Duration) Summary {
	s := Summary{
		Total:      len(results),
		DurationMs: duration.Milliseconds(),
	}
	for _, r := range results {
		if r.OK {
			s.Up++
		} else {
			s.Down++
		}
	}
	return s
}
