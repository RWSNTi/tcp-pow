package services

import (
	"encoding/json"
	"log"
	"pow-client/internal/models"
)

const (
	successResultKey = "success"
	failResultKet    = "fail"

	successResult = "You've succeeded solving the problem, your wise quote is:"
	failResult    = "You've failed, try again!"
)

var resultsMap = map[string]string{
	successResultKey: successResult,
	failResultKet:    failResult,
}

type challengeService struct {
	generateBlockService GenerateBlockService
}

type ChallengeService interface {
	SolveTheProblem(payload string) (responseBlock string, err error)
	PrintResult(payload string) error
}

func NewChallengeService(generateBlockService GenerateBlockService) ChallengeService {
	return &challengeService{
		generateBlockService: generateBlockService,
	}
}

// SolveTheProblem func to process solving the problem, accepts block object in JSON and responds with block object in JSON
func (s *challengeService) SolveTheProblem(payload string) (responseBlock string, err error) {
	var requestBlock models.Block
	err = json.Unmarshal([]byte(payload), &requestBlock)
	if err != nil {
		return
	}

	block, err := s.generateBlockService.GenerateBlock(requestBlock)
	if err != nil {
		return
	}

	log.Println("generated block with hash:", block.Hash)

	buf, err := json.Marshal(block)
	if err != nil {
		return
	}

	responseBlock = string(buf)
	return
}

// PrintResult is used to log and print result of work. Uses map to print predefined parts of result
func (s *challengeService) PrintResult(payload string) error {
	var result models.Result
	err := json.Unmarshal([]byte(payload), &result)
	if err != nil {
		return err
	}

	log.Println(resultsMap[result.Key], result.Quote)

	return nil
}
