package interfaces

import (
	"context"
	"telegram-quotes-bot/internal/entities"
)

// QuoteAPI определяет методы для работы с API цитат.
type QuoteAPI interface {
	GetRandomQuote(ctx context.Context) (*entities.Quote, error)
}

// Translator определяет методы для перевода текста.
type Translator interface {
	Translate(ctx context.Context, text, targetLang string) (string, error)
}

// TelegramSender определяет методы для отправки сообщений в Telegram.
type TelegramSender interface {
	SendMessage(ctx context.Context, message string) error
}

// CronScheduler определяет методы для планирования задач.
type CronScheduler interface {
	Start()
	AddJob(spec string, job func())
}
