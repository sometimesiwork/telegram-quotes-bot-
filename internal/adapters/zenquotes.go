package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"telegram-quotes-bot/internal/entities"
)

type ZenQuotesAPI struct{}

func NewZenQuotesAPI() *ZenQuotesAPI {
	return &ZenQuotesAPI{}
}

func (z *ZenQuotesAPI) GetRandomQuote(ctx context.Context) (*entities.Quote, error) {
	resp, err := http.Get("https://zenquotes.io/api/random")
	if err != nil {
		return nil, errors.New("ошибка запроса к API")
	}
	defer resp.Body.Close()

	var quotes []struct {
		Quote  string `json:"q"`
		Author string `json:"a"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
		return nil, errors.New("ошибка декодирования JSON")
	}

	if len(quotes) == 0 {
		return nil, errors.New("получен пустой список цитат")
	}

	return &entities.Quote{
		Text:   quotes[0].Quote,
		Author: quotes[0].Author,
	}, nil
}
