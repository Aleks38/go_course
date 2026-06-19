package api

type createChecksRequest struct {
	URLs    []string             `json:"urls"`
	Options *createChecksOptions `json:"options,omitempty"`
}
