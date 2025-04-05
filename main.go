package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log/slog"
	"os"
	"strconv"
	"telegram-quotes-bot/internal/adapters"
	"telegram-quotes-bot/internal/usecases"
)

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger
}

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Ошибка загрузки .env файла", "error", err)
		os.Exit(1)
	}

	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		slog.Error("Ошибка преобразования CHAT_ID в int64", "error", err)
		os.Exit(1)
	}

	logger := setupLogger()

	quoteAPI := adapters.NewZenQuotesAPI()
	translator := adapters.NewMyMemoryTranslator()
	telegramAdapter, err := adapters.NewTelegramAdapter(botToken, chatID)
	if err != nil {
		logger.Error("Не удалось инициализировать TelegramAdapter", "error", err)
		os.Exit(1)
	}

	fetchQuoteService := usecases.NewFetchQuoteService(quoteAPI)
	translateService := usecases.NewTranslateService(translator)
	sendQuoteService := usecases.NewSendQuoteService(telegramAdapter)

	c := cron.New()
	defer c.Stop()

	c.AddFunc("*/30 * * * *", func() {
		ctx := context.Background()

		quote, err := fetchQuoteService.FetchQuote(ctx)
		if err != nil {
			logger.Error("Ошибка получения цитаты", "error", err)
			return
		}

		translatedText, err := translateService.Translate(ctx, quote.Text)
		if err != nil {
			logger.Error("Ошибка перевода цитаты", "error", err)
		} else {
			quote.Text = translatedText
		}

		if err := sendQuoteService.SendQuote(ctx, quote); err != nil {
			logger.Error("Ошибка отправки цитаты", "error", err)
		} else {
			logger.Info("Цитата успешно отправлена", "quote", quote.Text)
		}
	})

	c.Start()
	logger.Info("Планировщик запущен. Ожидание задач.")

	select {}
}
