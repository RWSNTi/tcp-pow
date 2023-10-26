package services

import (
	"crypto/sha256"
	"encoding/hex"
	"pow-server/internal/models"
	"strconv"
	"strings"
	"sync"
	"time"
)

type blockchainService struct {
	mux         sync.Mutex
	targetValue string
	difficulty  int
	blockChain  []models.Block
}

type BlockchainService interface {
	CreateBlockchain()
	AddNewBlock(newBlock models.Block) bool
	GetLastBlock() models.Block
}

func NewBlockchainService(difficulty int, targetValue string) BlockchainService {
	return &blockchainService{
		difficulty:  difficulty,
		targetValue: targetValue,
	}
}

// CreateBlockchain creates a genesis block and adds it to blockchain
func (s *blockchainService) CreateBlockchain() {
	genesisBlock := models.Block{
		Index:      0,
		Timestamp:  time.Now().String(),
		Data:       0,
		Difficulty: s.difficulty,
		Target:     s.targetValue,
	}
	genesisBlock.Hash = s.calculateHash(genesisBlock)

	s.mux.Lock()
	s.blockChain = append(s.blockChain, genesisBlock)
	s.mux.Unlock()

	return
}

// GetLastBlock gets last block form blockchain
func (s *blockchainService) GetLastBlock() models.Block {
	return s.blockChain[len(s.blockChain)-1]
}

// AddNewBlock adds new block to blockchain after it's validation
func (s *blockchainService) AddNewBlock(newBlock models.Block) bool {
	if !s.validateBlock(newBlock) {
		return false
	}

	s.mux.Lock()
	s.blockChain = append(s.blockChain, newBlock)
	s.mux.Unlock()

	return true
}

// Validates block if the new block's fields and hash are appropriate
func (s *blockchainService) validateBlock(newBlock models.Block) bool {
	i := len(s.blockChain) - 1

	switch {
	case s.blockChain[i].Index+1 != newBlock.Index:
		return false
	case s.blockChain[i].Hash != newBlock.PrevHash:
		return false
	case newBlock.Target != s.targetValue:
		return false
	case s.calculateHash(newBlock) != newBlock.Hash:
		return false
	case !s.isHashValid(s.targetValue, newBlock.Hash, s.difficulty):
		return false
	}

	return true
}

// calculates hash for given block
func (s *blockchainService) calculateHash(block models.Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Data) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

// check if hash suits server requirements
func (s *blockchainService) isHashValid(target, hash string, difficulty int) bool {
	prefix := strings.Repeat(target, difficulty)
	return strings.HasPrefix(hash, prefix)
}
