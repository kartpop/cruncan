package model

type OneRequest struct {
	UserID string `json:"user_id"`
	Prompt string `json:"prompt"`
	Data   []Item `json:"data"`
}

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
