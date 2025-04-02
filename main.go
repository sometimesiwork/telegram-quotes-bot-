package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
// Замените значения на свои данные:
// - botToken – токен, полученный через BotFather
// - chatID   – идентификатор канала, например: -1001234567890
const (
	botToken = "8160500562:AAFi9TWrsZvltejKjXPI4vpzzXf59MmDwpY" // замените на ваш токен
	chatID   = -1002526755108                                   // замените на идентификатор канала (например, -1001234567890)
)

// zenQuoteURL – URL, к которому делается запрос для получения цитаты.
// Это значение можно переопределить в тестах.
var zenQuoteURL = "https://zenquotes.io/api/random"

// httpGet – функция для выполнения HTTP-запросов (по умолчанию http.Get).
// Позволяет переопределять её для тестирования.
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

	// Добавляем задачу – отправка цитаты каждый день в 9:00.
	_, err = c.AddFunc("0 9 * * *", func() {
		quote, err := fetchQuote()
		if err != nil {
			log.Printf("Ошибка получения цитаты: %v", err)
			return
		}

		msg := tgbotapi.NewMessage(chatID, quote)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		} else {
			log.Printf("Отправлена цитата: %s", quote)
		}
	})
	if err != nil {
		log.Fatalf("Ошибка добавления задачи в cron: %v", err)
	}

	// Запуск планировщика.
	c.Start()
	log.Println("Планировщик запущен. Ожидание задач.")

	// Запускаем тестовую отправку сразу (опционально, для проверки работы).
	// Если не требуется отправлять сразу, можно удалить этот блок.
	go func() {
		time.Sleep(2 * time.Second)
		quote, err := fetchQuote()
		if err != nil {
			log.Printf("Ошибка получения тестовой цитаты: %v", err)
			return
		}
		msg := tgbotapi.NewMessage(chatID, quote)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки тестового сообщения: %v", err)
		} else {
			log.Printf("Тестовая цитата успешно отправлена: %s", quote)
		}
	}()

	// Приложение работает бесконечно.
	select {}
}
