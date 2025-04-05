package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"log/slog"
	"os"
	"telegram-quotes-bot/internal/adapters"
	"telegram-quotes-bot/internal/config"
	"telegram-quotes-bot/internal/usecases"
)

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger
}

func main() {
	// Настройка логгера
	logger := setupLogger()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Error("Ошибка загрузки конфигурации", "error", err)
		os.Exit(1)
	}

	// Инициализация адаптеров
	quoteAPI := adapters.NewZenQuotesAPI()
	translator := adapters.NewMyMemoryTranslator()
	telegramAdapter, err := adapters.NewTelegramAdapter(cfg.BotToken, cfg.ChatID)
	if err != nil {
		logger.Error("Не удалось инициализировать TelegramAdapter", "error", err)
		os.Exit(1)
	}

	// Инициализация сервисов
	fetchQuoteService := usecases.NewFetchQuoteService(quoteAPI)
	translateService := usecases.NewTranslateService(translator)
	sendQuoteService := usecases.NewSendQuoteService(telegramAdapter)

	// Планировщик Cron
	c := cron.New()
	defer c.Stop()

	// Задача отправки цитат
	c.AddFunc("0 4,8,14,18 * * *", func() {
		ctx := context.Background()

		// Получение цитаты
		quote, err := fetchQuoteService.FetchQuote(ctx)
		if err != nil {
			logger.Error("Ошибка получения цитаты", "error", err)
			return
		}

		// Перевод цитаты
		translatedText, err := translateService.Translate(ctx, quote.Text)
		if err != nil {
			logger.Error("Ошибка перевода цитаты", "error", err)
		} else {
			quote.Text = translatedText
		}

		// Отправка цитаты
		if err := sendQuoteService.SendQuote(ctx, quote); err != nil {
			logger.Error("Ошибка отправки цитаты", "error", err)
		} else {
			logger.Info("Цитата успешно отправлена", "quote", quote.Text)
		}
	})

	// Запуск планировщика
	c.Start()
	logger.Info("Планировщик запущен. Ожидание задач.")

	// Бесконечный цикл для работы программы
	select {}
}
