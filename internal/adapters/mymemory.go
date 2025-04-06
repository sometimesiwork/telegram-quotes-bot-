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

// Translate выполняет перевод текста и имени автора через API MyMemory.
// Принимает контекст, исходный текст, имя автора и целевой язык (например, "ru").
// Возвращает переведённый текст, переведённое имя автора или ошибку, если запрос не удался.
func (t *MyMemoryTranslator) Translate(ctx context.Context, text, author, targetLang string) (string, string, error) {
	u := "https://api.mymemory.translated.net/get"
	params := url.Values{}

	// Перевод текста
	params.Set("q", text)
	params.Set("langpair", "en|"+targetLang)

	resp, err := http.Get(u + "?" + params.Encode())
	if err != nil {
		return "", "", errors.New("ошибка запроса к API при переводе текста")
	}
	defer resp.Body.Close()

	var result struct {
		ResponseData struct {
			TranslatedText string `json:"translatedText"`
		} `json:"responseData"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", errors.New("ошибка декодирования JSON при переводе текста")
	}

	if result.ResponseData.TranslatedText == "" {
		return "", "", errors.New("пустой ответ от MyMemory при переводе текста")
	}

	translatedText := result.ResponseData.TranslatedText

	// Перевод автора
	params.Set("q", author)

	resp, err = http.Get(u + "?" + params.Encode())
	if err != nil {
		return "", "", errors.New("ошибка запроса к API при переводе автора")
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", errors.New("ошибка декодирования JSON при переводе автора")
	}

	if result.ResponseData.TranslatedText == "" {
		return "", "", errors.New("пустой ответ от MyMemory при переводе автора")
	}

	translatedAuthor := result.ResponseData.TranslatedText

	return translatedText, translatedAuthor, nil
}
