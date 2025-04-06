package usecases

import (
	"context"
	"errors"

	"telegram-quotes-bot/internal/interfaces"
)

// TranslateService предоставляет методы для перевода текста и имени автора.
type TranslateService struct {
	translator interfaces.Translator // Интерфейс для выполнения перевода
}

// NewTranslateService создаёт новый экземпляр TranslateService.
// Принимает реализацию интерфейса Translator для выполнения перевода.
func NewTranslateService(translator interfaces.Translator) *TranslateService {
	return &TranslateService{translator: translator}
}

// Translate выполняет перевод текста и имени автора на русский язык.
// Принимает контекст, исходный текст и имя автора.
// Возвращает переведённый текст, переведённое имя автора или ошибку, если перевод не удался.
func (s *TranslateService) Translate(ctx context.Context, text, author string) (string, string, error) {
	// Вызываем метод Translate у переданного адаптера для перевода текста и автора
	translatedText, translatedAuthor, err := s.translator.Translate(ctx, text, author, "ru")
	if err != nil {
		// Если произошла ошибка при переводе, возвращаем пустые строки и сообщение об ошибке
		return "", "", errors.New("не удалось перевести текст или автора")
	}
	// Возвращаем переведённый текст и имя автора
	return translatedText, translatedAuthor, nil
}
