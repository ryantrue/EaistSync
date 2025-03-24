// commands.go
package telegrambot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartCommandListener запускает слушатель входящих обновлений и обрабатывает команды,
// вызывая callback-функцию для каждого сообщения с командой.
func (tb *TelegramBot) StartCommandListener(ctx context.Context, handler func(ctx context.Context, command string, args string, message tgbotapi.Message)) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := tb.botAPI.GetUpdatesChan(updateConfig)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				if update.Message == nil {
					continue
				}
				if update.Message.IsCommand() {
					cmd := update.Message.Command()
					args := update.Message.CommandArguments()
					handler(ctx, cmd, args, *update.Message)
				}
			}
		}
	}()
}
