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

func (s *TranslateService) Translate(ctx context.Context, text, author string) (string, string, error) {
	translatedText, translatedAuthor, err := s.translator.Translate(ctx, text, author, "ru")
	if err != nil {
		return "", "", errors.New("не удалось перевести текст или автора")
	}
	return translatedText, translatedAuthor, nil
}
