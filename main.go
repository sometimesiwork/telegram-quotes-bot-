package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	url2 "net/url"
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

// fetchQuote делает HTTP-запрос к API ZenQuotes и возвращает цитату в формате "Цитата – Автор".
func fetchQuote() (string, error) {
	resp, err := httpGet(zenQuoteURL)
	if err != nil {
		return "", fmt.Errorf("ошибка запроса к API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неожиданный HTTP статус: %d", resp.StatusCode)
	}

	// Распарсить JSON-ответ. Ответ представляет собой массив цитат.
	var quotes []Quote
	if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
		return "", fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	if len(quotes) == 0 {
		return "", fmt.Errorf("получен пустой список цитат")
	}

	// Форматируем ответ в виде: "Цитата – Автор"
	result := fmt.Sprintf("%s – %s", quotes[0].Quote, quotes[0].Author)
	return result, nil
}

// translateToRussian выполняет перевод текста на русский язык через MyMemory API.
func translateToRussian(text string) (string, error) {
	chunks := splitText(text)
	var translatedChunks []string

	for _, chunk := range chunks {
		translatedChunk, err := translateChunk(chunk)
		if err != nil {
			return "", err
		}
		translatedChunks = append(translatedChunks, translatedChunk)
	}

	return strings.Join(translatedChunks, " "), nil
}

// splitText разбивает текст на части длиной до maxTextLength символов.
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

// translateChunk переводит одну часть текста.
func translateChunk(text string) (string, error) {
	url := "https://api.mymemory.translated.net/get"

	// Кодируем текст для безопасного использования в URL
	encodedText := url2.QueryEscape(text)
	params := fmt.Sprintf("?q=%s&langpair=en|ru", encodedText)

	log.Printf("Выполняется запрос к MyMemory API: %s%s", url, params)

	resp, err := http.Get(url + params)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неожиданный HTTP статус: %d", resp.StatusCode)
	}

	var result struct {
		ResponseData struct {
			TranslatedText string `json:"translatedText"`
		} `json:"responseData"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	if result.ResponseData.TranslatedText == "" {
		return "", fmt.Errorf("пустой ответ от MyMemory")
	}

	return result.ResponseData.TranslatedText, nil
}

func sendQuote(bot *tgbotapi.BotAPI, chatID int64) {
	quote, err := fetchQuote()
	if err != nil {
		log.Printf("Ошибка получения цитаты: %v", err)
		return
	}

	// Переводим цитату на русский язык
	translatedQuote, err := translateToRussian(quote)
	if err != nil {
		log.Printf("Ошибка перевода цитаты: %v", err)
		translatedQuote = quote // Возвращаем оригинальную цитату, если перевод не удался
	}

	// Отправляем цитату в Telegram-канал
	msg := tgbotapi.NewMessage(chatID, translatedQuote)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	} else {
		log.Printf("Отправлена цитата: %s", translatedQuote)
	}
}

func main() {
	// Инициализация Telegram-бота.
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panicf("Не удалось инициализировать бота: %v", err)
	}
	bot.Debug = true
	log.Printf("Бот запущен под именем: %s", bot.Self.UserName)

	// Инициализация планировщика cron.
	c := cron.New()

	// Добавляем задачи – отправка цитаты три раза в день.
	times := []string{"0 9 * * *", "0 15 * * *", "0 21 * * *"} // 9:00, 15:00, 21:00
	for _, cronTime := range times {
		_, err := c.AddFunc(cronTime, func() {
			sendQuote(bot, chatID)
		})
		if err != nil {
			log.Fatalf("Ошибка добавления задачи в cron: %v", err)
		}
	}

	// Запуск планировщика.
	c.Start()
	log.Println("Планировщик запущен. Ожидание задач.")

	// Тестовая отправка цитаты
	go func() {
		time.Sleep(5 * time.Second) // Ждём 5 секунд перед отправкой тестовой цитаты
		sendQuote(bot, chatID)
	}()

	// Приложение работает бесконечно.
	select {}
}
