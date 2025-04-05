package usecases

import (
	"context"
	"errors"

	"telegram-quotes-bot/internal/entities"
	"telegram-quotes-bot/internal/interfaces"
)

type FetchQuoteService struct {
	api interfaces.QuoteAPI
}

func NewFetchQuoteService(api interfaces.QuoteAPI) *FetchQuoteService {
	return &FetchQuoteService{api: api}
}

func (s *FetchQuoteService) FetchQuote(ctx context.Context) (*entities.Quote, error) {
	quote, err := s.api.GetRandomQuote(ctx)
	if err != nil {
		return nil, errors.New("не удалось получить цитату")
	}
	return quote, nil
}
