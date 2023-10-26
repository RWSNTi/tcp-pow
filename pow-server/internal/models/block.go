package models

// Block model is used in Blockchain, also it contains challenge parameters: target, difficulty
type Block struct {
	Index      int    `json:"index"`
	Timestamp  string `json:"timestamp"`
	Data       int    `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prev_hash"`
	Difficulty int    `json:"difficulty"`
	Target     string `json:"target"`
	Nonce      string `json:"nonce"`
}
