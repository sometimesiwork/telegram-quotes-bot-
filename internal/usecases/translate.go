package usecases

import (
	"context"
	"errors"

	"telegram-quotes-bot/internal/interfaces"
)

type TranslateService struct {
	translator interfaces.Translator
}

// NewTranslateService создаёт новый экземпляр TranslateService.
// Принимает интерфейс Translator для выполнения перевода текста.
func NewTranslateService(translator interfaces.Translator) *TranslateService {
	return &TranslateService{translator: translator}
}

// Translate выполняет перевод текста на русский язык.
// Возвращает переведённый текст или ошибку, если перевод не удался.
func (s *TranslateService) Translate(ctx context.Context, text string) (string, error) {
	translatedText, err := s.translator.Translate(ctx, text, "ru")
	if err != nil {
		return "", errors.New("не удалось перевести текст")
	}
	return translatedText, nil
}
