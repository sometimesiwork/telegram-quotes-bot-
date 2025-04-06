package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// MyMemoryTranslator реализует интерфейс Translator для перевода текста через API MyMemory.
type MyMemoryTranslator struct{}

// NewMyMemoryTranslator создаёт новый экземпляр MyMemoryTranslator.
func NewMyMemoryTranslator() *MyMemoryTranslator {
	return &MyMemoryTranslator{}
}

// Translate выполняет перевод текста и имени автора через API MyMemory.
// Принимает контекст, исходный текст, имя автора и целевой язык (например, "ru").
// Возвращает переведённый текст, переведённое имя автора или ошибку, если запрос не удался.
func (t *MyMemoryTranslator) Translate(ctx context.Context, text, author, targetLang string) (string, string, error) {
	u := "https://api.mymemory.translated.net/get" // URL для запроса к API MyMemory
	params := url.Values{}

	// Перевод текста
	params.Set("q", text)                    // Устанавливаем текст для перевода
	params.Set("langpair", "en|"+targetLang) // Указываем пару языков для перевода (например, "en|ru")

	resp, err := http.Get(u + "?" + params.Encode()) // Выполняем GET-запрос к API
	if err != nil {
		return "", "", errors.New("ошибка запроса к API при переводе текста")
	}
	defer resp.Body.Close() // Закрываем тело ответа после завершения работы

	// Декодируем JSON-ответ от API
	var result struct {
		ResponseData struct {
			TranslatedText string `json:"translatedText"` // Поле с переведённым текстом
		} `json:"responseData"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", errors.New("ошибка декодирования JSON при переводе текста")
	}

	// Проверяем, что ответ содержит переведённый текст
	if result.ResponseData.TranslatedText == "" {
		return "", "", errors.New("пустой ответ от MyMemory при переводе текста")
	}

	translatedText := result.ResponseData.TranslatedText // Сохраняем переведённый текст

	// Перевод имени автора
	params.Set("q", author) // Устанавливаем имя автора для перевода

	resp, err = http.Get(u + "?" + params.Encode()) // Выполняем GET-запрос к API
	if err != nil {
		return "", "", errors.New("ошибка запроса к API при переводе автора")
	}
	defer resp.Body.Close() // Закрываем тело ответа после завершения работы

	// Декодируем JSON-ответ от API
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", errors.New("ошибка декодирования JSON при переводе автора")
	}

	// Проверяем, что ответ содержит переведённое имя автора
	if result.ResponseData.TranslatedText == "" {
		return "", "", errors.New("пустой ответ от MyMemory при переводе автора")
	}

	translatedAuthor := result.ResponseData.TranslatedText // Сохраняем переведённое имя автора

	// Возвращаем переведённые текст и имя автора
	return translatedText, translatedAuthor, nil
}
