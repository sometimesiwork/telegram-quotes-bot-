package usecases

import (
	"context"
	"errors"

	"telegram-quotes-bot/internal/interfaces"
)

type TranslateService struct {
	translator interfaces.Translator
}

func NewTranslateService(translator interfaces.Translator) *TranslateService {
	return &TranslateService{translator: translator}
}

func (s *TranslateService) Translate(ctx context.Context, text string) (string, error) {
	translatedText, err := s.translator.Translate(ctx, text, "ru")
	if err != nil {
		return "", errors.New("не удалось перевести текст")
	}
	return translatedText, nil
}
