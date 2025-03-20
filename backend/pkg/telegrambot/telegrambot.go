package telegrambot

import (
	"context"
	"fmt"
	"log"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramBot представляет Telegram-бота для отправки сообщений.
type TelegramBot struct {
	Bot        *tgbotapi.BotAPI
	ChatID     int64         // Идентификатор чата для отправки сообщений
	MaxRetries int           // Максимальное количество попыток отправки сообщения при ошибке
	RetryDelay time.Duration // Задержка между попытками отправки сообщения
}

// NewTelegramBot создает и возвращает новый экземпляр TelegramBot.
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

// sleepWithContext осуществляет ожидание с учетом отмены context.
func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// sendMessageChunk отправляет одну часть сообщения с повторными попытками.
func (tb *TelegramBot) sendMessageChunk(ctx context.Context, message string) error {
	msg := tgbotapi.NewMessage(tb.ChatID, message)
	var lastErr error

	for attempt := 0; attempt <= tb.MaxRetries; attempt++ {
		// Проверка на отмену context.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, err := tb.Bot.Send(msg)
		if err == nil {
			return nil
		}

		lastErr = err
		log.Printf("Attempt %d: failed to send message: %v", attempt+1, err)

		// Если это не последняя попытка, ждем перед повторной отправкой.
		if attempt < tb.MaxRetries {
			if err := sleepWithContext(ctx, tb.RetryDelay); err != nil {
				return err
			}
		}
	}

	return fmt.Errorf("failed to send message after %d attempts: %w", tb.MaxRetries+1, lastErr)
}

// splitMessage разбивает сообщение на части длиной не более maxLength символов.
func splitMessage(message string, maxLength int) []string {
	var parts []string
	runes := []rune(message)
	for len(runes) > 0 {
		if len(runes) <= maxLength {
			parts = append(parts, string(runes))
			break
		}
		parts = append(parts, string(runes[:maxLength]))
		runes = runes[maxLength:]
	}
	return parts
}

// SendMessage отправляет сообщение через Telegram-бота.
// Если сообщение превышает максимальную длину (4096 символов), оно разбивается на части и отправляется последовательно.
func (tb *TelegramBot) SendMessage(ctx context.Context, message string) error {
	const maxLength = 4096

	if utf8.RuneCountInString(message) <= maxLength {
		return tb.sendMessageChunk(ctx, message)
	}

	parts := splitMessage(message, maxLength)
	for i, part := range parts {
		log.Printf("Отправка части сообщения %d из %d", i+1, len(parts))
		if err := tb.sendMessageChunk(ctx, part); err != nil {
			return fmt.Errorf("failed to send message part %d: %w", i+1, err)
		}
	}
	return nil
}

// Notify реализует интерфейс Notifier.
func (tb *TelegramBot) Notify(ctx context.Context, message string) error {
	return tb.SendMessage(ctx, message)
}
