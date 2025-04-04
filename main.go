package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

// Quote соответствует структуре ответа от ZenQuotes API.
type Quote struct {
	Quote  string `json:"q"` // текст цитаты
	Author string `json:"a"` // автор
	HTML   string `json:"h"` // HTML-версия (необязательно)
}

// Константы для токена бота и идентификатора канала.
const (
	botToken = "8160500562:AAFi9TWrsZvltejKjXPI4vpzzXf59MmDwpY" // замените на ваш токен
	chatID   = -1002526755108                                   // замените на идентификатор канала
)

// zenQuoteURL – URL, к которому делается запрос для получения цитаты.
var zenQuoteURL = "https://zenquotes.io/api/random"

// httpGet – функция для выполнения HTTP-запросов (по умолчанию http.Get).
var httpGet = http.Get

// Настройка логгера slog
func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return logger
}

// fetchQuote делает HTTP-запрос к API ZenQuotes и возвращает цитату в формате "Цитата – Автор".
func fetchQuote(logger *slog.Logger) (string, error) {
	resp, err := httpGet(zenQuoteURL)
	if err != nil {
		logger.Error("Ошибка запроса к API", "error", err)
		return "", fmt.Errorf("ошибка запроса к API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Неожиданный HTTP статус", "status_code", resp.StatusCode)
		return "", fmt.Errorf("неожиданный HTTP статус: %d", resp.StatusCode)
	}

	var quotes []Quote
	if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
		logger.Error("Ошибка декодирования JSON", "error", err)
		return "", fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	if len(quotes) == 0 {
		logger.Error("Получен пустой список цитат")
		return "", fmt.Errorf("получен пустой список цитат")
	}

	result := fmt.Sprintf("%s – %s", quotes[0].Quote, quotes[0].Author)
	logger.Info("Цитата успешно получена", "quote", result)
	return result, nil
}

// translateToRussian выполняет перевод текста на русский язык через MyMemory API.
func translateToRussian(text string, logger *slog.Logger) (string, error) {
	chunks := splitText(text)
	var translatedChunks []string

	for _, chunk := range chunks {
		translatedChunk, err := translateChunk(chunk, logger)
		if err != nil {
			logger.Error("Ошибка перевода части текста", "chunk", chunk, "error", err)
			return "", err
		}
		translatedChunks = append(translatedChunks, translatedChunk)
	}

	translatedText := strings.Join(translatedChunks, " ")
	logger.Info("Текст успешно переведён", "original", text, "translated", translatedText)
	return translatedText, nil
}

const maxTextLength = 500

func splitText(text string) []string {
	var chunks []string
	for len(text) > maxTextLength {
		chunks = append(chunks, text[:maxTextLength])
		text = text[maxTextLength:]
	}
	chunks = append(chunks, text)
	return chunks
}

func translateChunk(text string, logger *slog.Logger) (string, error) {
	url := "https://api.mymemory.translated.net/get"
	encodedText := url2.QueryEscape(text)
	params := fmt.Sprintf("?q=%s&langpair=en|ru", encodedText)

	logger.Info("Выполняется запрос к MyMemory API", "url", url+params)

	resp, err := http.Get(url + params)
	if err != nil {
		logger.Error("Ошибка при выполнении HTTP-запроса", "url", url+params, "error", err)
		return "", fmt.Errorf("ошибка при выполнении HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Неожиданный HTTP статус", "url", url+params, "status_code", resp.StatusCode)
		return "", fmt.Errorf("неожиданный HTTP статус: %d", resp.StatusCode)
	}

	var result struct {
		ResponseData struct {
			TranslatedText string `json:"translatedText"`
		} `json:"responseData"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("Ошибка декодирования JSON", "url", url+params, "error", err)
		return "", fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	if result.ResponseData.TranslatedText == "" {
		logger.Error("Пустой ответ от MyMemory", "url", url+params)
		return "", fmt.Errorf("пустой ответ от MyMemory")
	}

	return result.ResponseData.TranslatedText, nil
}

func sendQuote(bot *tgbotapi.BotAPI, chatID int64, logger *slog.Logger) {
	logger.Info("Задача отправки цитаты запущена")
	quote, err := fetchQuote(logger)
	if err != nil {
		logger.Error("Ошибка получения цитаты", "error", err)
		return
	}

	translatedQuote, err := translateToRussian(quote, logger)
	if err != nil {
		logger.Error("Ошибка перевода цитаты", "error", err)
		translatedQuote = quote // Если перевод не удался, используем оригинальную цитату
	}

	// Разделяем цитату на текст и автора
	parts := strings.Split(translatedQuote, " – ")
	if len(parts) != 2 {
		logger.Error("Некорректный формат цитаты", "quote", translatedQuote)
		return
	}
	quoteText := parts[0]
	quoteAuthor := parts[1]

	// Формируем сообщение с эмодзи
	formattedMessage := fmt.Sprintf(
		"📖 %s\n\n— %s ✍️",
		quoteText,
		quoteAuthor,
	)

	msg := tgbotapi.NewMessage(chatID, formattedMessage)
	if _, err := bot.Send(msg); err != nil {
		logger.Error("Ошибка отправки сообщения", "error", err)
	} else {
		logger.Info("Цитата успешно отправлена", "quote", translatedQuote)
	}
}

func main() {
	// Настройка логгера
	logger := setupLogger()

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.Error("Не удалось инициализировать бота", "error", err)
		os.Exit(1)
	}
	bot.Debug = true
	logger.Info("Бот запущен", "bot_name", bot.Self.UserName)

	// Инициализация планировщика Cron
	c := cron.New()
	defer c.Stop()

	// Задачи отправки цитат три раза в день (время в UTC)
	times := []string{"0 3 * * *", "0 9 * * *", "0 15 * * *"} // 6:00, 12:00, 18:00 МСК
	for _, cronTime := range times {
		_, err := c.AddFunc(cronTime, func() {
			sendQuote(bot, chatID, logger)
		})
		if err != nil {
			logger.Error("Ошибка добавления задачи в cron", "error", err)
			os.Exit(1)
		}
	}

	c.Start()
	logger.Info("Планировщик запущен. Ожидание задач.")

	// Тестовый запрос через 5 секунд после запуска
	go func() {
		time.Sleep(5 * time.Second) // Ждём 5 секунд
		logger.Info("Выполняется тестовый запрос...")
		sendQuote(bot, chatID, logger)
	}()

	select {}
}
