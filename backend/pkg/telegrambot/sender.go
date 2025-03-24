// sender.go
package telegrambot

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// telegramSender реализует интерфейс BotSender.
type telegramSender struct {
	bot        *tgbotapi.BotAPI
	chatID     int64
	maxRetries int
	retryDelay time.Duration
}

// newTelegramSender создаёт нового отправщика.
func newTelegramSender(bot *tgbotapi.BotAPI, chatID int64, maxRetries int, retryDelay time.Duration) *telegramSender {
	return &telegramSender{
		bot:        bot,
		chatID:     chatID,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// ChatID возвращает идентификатор чата.
func (ts *telegramSender) ChatID() int64 {
	return ts.chatID
}

// SendMessage отправляет текстовое сообщение, используя механизм повторных попыток.
func (ts *telegramSender) SendMessage(ctx context.Context, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	operation := func() error {
		_, err := ts.bot.Send(msg)
		return err
	}
	return retryOperation(ctx, operation, ts.maxRetries, ts.retryDelay)
}

// SendDocument отправляет документ, используя механизм повторных попыток.
func (ts *telegramSender) SendDocument(ctx context.Context, chatID int64, fileName string, content []byte) error {
	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: content,
	})
	operation := func() error {
		_, err := ts.bot.Send(doc)
		return err
	}
	return retryOperation(ctx, operation, ts.maxRetries, ts.retryDelay)
}
