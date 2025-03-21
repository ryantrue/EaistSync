package telegrambot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramBot представляет Telegram-бота для отправки сообщений и документов.
type TelegramBot struct {
	Bot        *tgbotapi.BotAPI
	ChatID     int64         // Идентификатор чата для отправки сообщений/документов
	MaxRetries int           // Максимальное количество попыток отправки
	RetryDelay time.Duration // Задержка между попытками
}

// NewTelegramBot создаёт и возвращает новый экземпляр TelegramBot.
func NewTelegramBot(token string, chatID int64, maxRetries int, retryDelay time.Duration) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error initializing Telegram bot: %w", err)
	}
	return &TelegramBot{
		Bot:        bot,
		ChatID:     chatID,
		MaxRetries: maxRetries,
		RetryDelay: retryDelay,
	}, nil
}

// retryOperation выполняет указанную операцию с повторными попытками.
// Если операция завершается успешно (err == nil) – возвращается nil.
// Если контекст отменён или исчерпаны попытки, возвращается последняя ошибка.
func retryOperation(ctx context.Context, operation func() error, maxRetries int, retryDelay time.Duration) error {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Проверка контекста.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := operation(); err != nil {
			lastErr = err
			log.Printf("Attempt %d failed: %v", attempt+1, err)
			// Ждём либо время retryDelay, либо отмены контекста.
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay):
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries+1, lastErr)
}

// SendJSONDocument сериализует произвольную структуру в JSON и отправляет её как документ.
// Данные сериализуются в память и передаются через tgbotapi.FileBytes.
func (tb *TelegramBot) SendJSONDocument(ctx context.Context, document interface{}) error {
	// Сериализуем документ в JSON с отступами для читаемости.
	jsonData, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal document into JSON: %w", err)
	}

	// Создаём документ для отправки, используя данные из памяти.
	doc := tgbotapi.NewDocument(tb.ChatID, tgbotapi.FileBytes{
		Name:  "document.json",
		Bytes: jsonData,
	})

	// Операция отправки документа.
	operation := func() error {
		_, err := tb.Bot.Send(doc)
		return err
	}

	return retryOperation(ctx, operation, tb.MaxRetries, tb.RetryDelay)
}

// Notify отправляет короткое текстовое сообщение (например, для уведомлений или ошибок).
func (tb *TelegramBot) Notify(ctx context.Context, message string) error {
	msg := tgbotapi.NewMessage(tb.ChatID, message)
	operation := func() error {
		_, err := tb.Bot.Send(msg)
		return err
	}
	return retryOperation(ctx, operation, tb.MaxRetries, tb.RetryDelay)
}
