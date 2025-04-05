package usecases

import (
	"context"
	"fmt"

	"telegram-quotes-bot/internal/entities"
	"telegram-quotes-bot/internal/interfaces"
)

type SendQuoteService struct {
	telegram interfaces.TelegramSender
}

func NewSendQuoteService(telegram interfaces.TelegramSender) *SendQuoteService {
	return &SendQuoteService{telegram: telegram}
}

func (s *SendQuoteService) SendQuote(ctx context.Context, quote *entities.Quote) error {
	message := fmt.Sprintf("📖 %s\n\n— %s ✍️", quote.Text, quote.Author)
	err := s.telegram.SendMessage(ctx, message)
	if err != nil {
		return fmt.Errorf("не удалось отправить сообщение: %w", err)
	}
	return nil
}
