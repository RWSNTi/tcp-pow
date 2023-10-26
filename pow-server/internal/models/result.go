package models

// Result model is used to parse result response from server
type Result struct {
	Key    string `json:"key"`
	Result string `json:"result"`
}
