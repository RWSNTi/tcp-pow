package models

// TCPExchange model is used for message exchange between server and client
type TCPExchange struct {
	Header  string `json:"header"`
	Payload string `json:"payload"`
}
