package models

// Quote model is used to parse message from Quote API server.
type Quote struct {
	Quote string `json:"q"`
}
