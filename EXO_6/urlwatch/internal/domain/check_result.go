package domain

type CheckResult struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code,omitempty"`
	OK         bool   `json:"ok"`
	LatencyMs  int64  `json:"latency_ms"`
	Error      string `json:"error,omitempty"`
}
