// telegrambot.go
package telegrambot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotSender определяет интерфейс отправки сообщений и документов.
type BotSender interface {
	// ChatID возвращает идентификатор чата.
	ChatID() int64
	// SendMessage отправляет текстовое сообщение.
	SendMessage(ctx context.Context, chatID int64, text string) error
	// SendDocument отправляет документ с указанным именем и содержимым.
	SendDocument(ctx context.Context, chatID int64, fileName string, content []byte) error
}

// TelegramBot – обёртка для отправки уведомлений через Telegram, теперь также включает возможность обработки команд.
type TelegramBot struct {
	sender BotSender
	botAPI *tgbotapi.BotAPI // для обработки входящих обновлений (команд)
}

// NewTelegramBot создаёт экземпляр TelegramBot, инициализируя внутреннего отправщика.
func NewTelegramBot(token string, chatID int64, maxRetries int, retryDelay time.Duration) (*TelegramBot, error) {
	bot, err := newBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error initializing Telegram bot: %w", err)
	}
	sender := newTelegramSender(bot, chatID, maxRetries, retryDelay)
	return &TelegramBot{
		sender: sender,
		botAPI: bot,
	}, nil
}

// Notify отправляет текстовое сообщение для уведомлений или ошибок.
func (tb *TelegramBot) Notify(ctx context.Context, message string) error {
	return tb.sender.SendMessage(ctx, tb.sender.ChatID(), message)
}

// SendJSONDocument сериализует структуру в JSON и отправляет её как документ
// с дефолтным именем файла "document.json".
func (tb *TelegramBot) SendJSONDocument(ctx context.Context, document interface{}) error {
	return tb.SendJSONDocumentWithName(ctx, document, "document.json")
}

// SendJSONDocumentWithName сериализует структуру в JSON и отправляет её как документ,
// используя указанное имя файла.
func (tb *TelegramBot) SendJSONDocumentWithName(ctx context.Context, document interface{}, fileName string) error {
	jsonData, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal document into JSON: %w", err)
	}
	return tb.sender.SendDocument(ctx, tb.sender.ChatID(), fileName, jsonData)
}
