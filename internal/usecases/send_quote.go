package usecases

import (
	"context"
	"fmt"

	"telegram-quotes-bot/internal/entities"
	"telegram-quotes-bot/internal/interfaces"
)

// SendQuoteService предоставляет методы для отправки цитат в Telegram-канал.
type SendQuoteService struct {
	telegram interfaces.TelegramSender // Интерфейс для отправки сообщений в Telegram
}

// NewSendQuoteService создаёт новый экземпляр SendQuoteService.
// Принимает интерфейс TelegramSender для отправки сообщений в Telegram.
func NewSendQuoteService(telegram interfaces.TelegramSender) *SendQuoteService {
	return &SendQuoteService{telegram: telegram}
}

// SendQuote отправляет цитату в Telegram-канал.
// Форматирует цитату в удобочитаемый вид и отправляет её через TelegramSender.
// Возвращает ошибку, если отправка не удалась.
func (s *SendQuoteService) SendQuote(ctx context.Context, quote *entities.Quote) error {
	// Форматируем цитату с эмодзи для текста и автора
	message := fmt.Sprintf("📖 *%s*\n\n— _%s_ ✍️", quote.Text, quote.Author)

	// Отправляем сформированное сообщение через TelegramSender
	err := s.telegram.SendMessage(ctx, message)
	if err != nil {
		// Если произошла ошибка при отправке, возвращаем её с описанием
		return fmt.Errorf("не удалось отправить сообщение: %w", err)
	}

	// Если всё прошло успешно, возвращаем nil
	return nil
}
