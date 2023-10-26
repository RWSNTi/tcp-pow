package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"pow-client/internal/models"
	"strconv"
	"strings"
	"time"
)

type generateBlockService struct {
	limit int
}

type GenerateBlockService interface {
	GenerateBlock(requestBlock models.Block) (newBlock models.Block, err error)
}

func NewGenerateBlockService(limit int) GenerateBlockService {
	return &generateBlockService{
		limit: limit,
	}
}

func(s *generateBlockService) GenerateBlock(requestBlock models.Block) (newBlock models.Block, err error) {

	newBlock.Index = requestBlock.Index + 1
	newBlock.Data = 0
	newBlock.PrevHash = requestBlock.Hash
	newBlock.Difficulty = requestBlock.Difficulty
	newBlock.Timestamp = time.Now().String()
	newBlock.Target = requestBlock.Target

	for i := 0; i < s.limit; i++ {
		hexValue := fmt.Sprintf("%x", i)
		newBlock.Nonce = hexValue

		hash := s.calculateHash(newBlock)
//		fmt.Printf("\r%s", hash) // Uncomment to see hash generation process
		if !s.isHashValid(requestBlock.Target, hash, requestBlock.Difficulty) {
			continue
		}

		newBlock.Hash = hash
//		fmt.Println("")  // Uncomment to see hash generation process
		return newBlock, nil
	}

	err = fmt.Errorf("couldn't create suitable hash in required number of iterations")
	return
}

func(s *generateBlockService) calculateHash(block models.Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Data) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}

func(s *generateBlockService) isHashValid(target, hash string, difficulty int) bool {
	prefix := strings.Repeat(target, difficulty)
	return strings.HasPrefix(hash, prefix)
}
