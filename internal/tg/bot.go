package tg

import (
	"fmt"
	"github.com/SerGGe43/alertilka/internal/bot"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbot.BotAPI
}

func NewBot(api *tgbot.BotAPI) *Bot {
	return &Bot{
		api: api,
	}
}

func (b *Bot) SendMenu(chatID int64) error {
	var mainMenuKeyboard = tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton(bot.Help)),
		tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton(bot.PriceByTicker)),
		tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton(bot.NewAlert)),
	)
	msg := tgbot.NewMessage(chatID, "Main menu")
	msg.ReplyMarkup = mainMenuKeyboard
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send menu: %w", err)
	}
	return nil
}
