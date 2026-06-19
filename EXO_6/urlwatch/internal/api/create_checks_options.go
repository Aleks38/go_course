package api

type createChecksOptions struct {
	Concurrency *int `json:"concurrency,omitempty"`
	TimeoutMs   *int `json:"timeout_ms,omitempty"`
}
