package domain

type Summary struct {
	Total      int   `json:"total"`
	Up         int   `json:"up"`
	Down       int   `json:"down"`
	DurationMs int64 `json:"duration_ms"`
}
