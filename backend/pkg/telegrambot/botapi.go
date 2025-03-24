// botapi.go
package telegrambot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// newBotAPI инициализирует и возвращает новый экземпляр BotAPI.
// При необходимости эту функцию можно заменить мок-реализацией для тестов.
func newBotAPI(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Telegram bot: %w", err)
	}
	return bot, nil
}
