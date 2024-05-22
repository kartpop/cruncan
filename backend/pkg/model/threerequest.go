package model

// ThreeRequest is a model for the request to /three
type ThreeRequest struct {
	ID         string     `json:"id"`
	Metadata   string     `json:"metadata"`
	OneRequest OneRequest `json:"one_request"`
}
