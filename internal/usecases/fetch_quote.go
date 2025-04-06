package usecases

import (
	"context"
	"errors"

	"telegram-quotes-bot/internal/entities"
	"telegram-quotes-bot/internal/interfaces"
)

// FetchQuoteService предоставляет методы для получения случайных цитат через API.
type FetchQuoteService struct {
	api interfaces.QuoteAPI // Интерфейс для взаимодействия с внешним API цитат
}

// NewFetchQuoteService создаёт новый экземпляр FetchQuoteService.
// Принимает реализацию интерфейса QuoteAPI для получения цитат.
func NewFetchQuoteService(api interfaces.QuoteAPI) *FetchQuoteService {
	return &FetchQuoteService{api: api}
}

// FetchQuote получает случайную цитату через API.
// Возвращает структуру Quote или ошибку, если не удалось получить цитату.
func (s *FetchQuoteService) FetchQuote(ctx context.Context) (*entities.Quote, error) {
	// Вызываем метод GetRandomQuote у переданного API для получения случайной цитаты
	quote, err := s.api.GetRandomQuote(ctx)
	if err != nil {
		// Если произошла ошибка при получении цитаты, возвращаем nil и сообщение об ошибке
		return nil, errors.New("не удалось получить цитату")
	}
	// Возвращаем полученную цитату
	return quote, nil
}
