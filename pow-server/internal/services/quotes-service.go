package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pow-server/internal/models"
)

type quotesService struct {
	apiUrl string
}

type QuotesService interface {
	GetQuote() (result string, err error)
}

func NewQuotesService(apiUrl string) QuotesService {
	return &quotesService{
		apiUrl: apiUrl,
	}
}

func(s *quotesService) GetQuote() (result string, err error) {
	resp, err := http.Get(s.apiUrl + "/random")
	if err != nil {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var quotes []models.Quote
	err = json.Unmarshal(body, &quotes)
	if err != nil {
		return
	}

	if len(quotes) == 0 {
		err = fmt.Errorf("no quotes received from source")
		return
	}

	return quotes[0].Quote, nil
}
