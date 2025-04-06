package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"telegram-quotes-bot/internal/entities"
)

// ZenQuotesAPI реализует интерфейс QuoteAPI для получения случайных цитат из ZenQuotes API.
type ZenQuotesAPI struct{}

// NewZenQuotesAPI создаёт новый экземпляр ZenQuotesAPI.
func NewZenQuotesAPI() *ZenQuotesAPI {
	return &ZenQuotesAPI{}
}

// GetRandomQuote получает случайную цитату из ZenQuotes API.
// Возвращает структуру Quote или ошибку, если запрос или декодирование не удались.
func (z *ZenQuotesAPI) GetRandomQuote(ctx context.Context) (*entities.Quote, error) {
	// Выполняем GET-запрос к ZenQuotes API для получения случайной цитаты
	resp, err := http.Get("https://zenquotes.io/api/random")
	if err != nil {
		// Если произошла ошибка при выполнении запроса, возвращаем её
		return nil, errors.New("ошибка запроса к API")
	}
	defer resp.Body.Close() // Закрываем тело ответа после завершения работы

	// Определяем структуру для декодирования JSON-ответа
	var quotes []struct {
		Quote  string `json:"q"` // Текст цитаты
		Author string `json:"a"` // Имя автора
	}

	// Декодируем JSON-ответ от API
	if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
		// Если произошла ошибка при декодировании JSON, возвращаем её
		return nil, errors.New("ошибка декодирования JSON")
	}

	// Проверяем, что ответ содержит хотя бы одну цитату
	if len(quotes) == 0 {
		return nil, errors.New("получен пустой список цитат")
	}

	// Возвращаем цитату, преобразованную в структуру Quote
	return &entities.Quote{
		Text:   quotes[0].Quote,  // Текст цитаты
		Author: quotes[0].Author, // Имя автора
	}, nil
}
