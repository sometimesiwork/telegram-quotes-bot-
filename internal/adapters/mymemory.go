package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type MyMemoryTranslator struct{}

// NewMyMemoryTranslator создаёт новый экземпляр MyMemoryTranslator.
func NewMyMemoryTranslator() *MyMemoryTranslator {
	return &MyMemoryTranslator{}
}

// Translate выполняет перевод текста через API MyMemory.
// Принимает контекст, исходный текст и целевой язык (например, "ru").
// Возвращает переведённый текст или ошибку, если запрос не удался.
func (t *MyMemoryTranslator) Translate(ctx context.Context, text, targetLang string) (string, error) {
	u := "https://api.mymemory.translated.net/get"
	params := url.Values{}
	params.Set("q", text)
	params.Set("langpair", "en|"+targetLang)

	resp, err := http.Get(u + "?" + params.Encode())
	if err != nil {
		return "", errors.New("ошибка запроса к API")
	}
	defer resp.Body.Close()

	var result struct {
		ResponseData struct {
			TranslatedText string `json:"translatedText"`
		} `json:"responseData"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", errors.New("ошибка декодирования JSON")
	}

	if result.ResponseData.TranslatedText == "" {
		return "", errors.New("пустой ответ от MyMemory")
	}

	return result.ResponseData.TranslatedText, nil
}
