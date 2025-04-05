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

// NewFetchQuoteService создаёт новый экземпляр FetchQuoteService.
// Принимает интерфейс QuoteAPI для получения цитат.
func NewFetchQuoteService(api interfaces.QuoteAPI) *FetchQuoteService {
	return &FetchQuoteService{api: api}
}

// FetchQuote получает случайную цитату через API.
// Возвращает структуру Quote или ошибку, если не удалось получить цитату.
func (s *FetchQuoteService) FetchQuote(ctx context.Context) (*entities.Quote, error) {
	quote, err := s.api.GetRandomQuote(ctx)
	if err != nil {
		return nil, errors.New("не удалось получить цитату")
	}
	return quote, nil
}
