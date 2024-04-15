package model

type OneRequest struct {
	UserID string `json:"user_id"`
	Prompt string `json:"prompt"`
	Data  Data `json:"data"`

}

type Data struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
