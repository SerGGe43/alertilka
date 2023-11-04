package tg

import (
	"fmt"
	"github.com/SerGGe43/alertilka/internal/bot"
	"github.com/SerGGe43/alertilka/internal/domain"
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

func (b *Bot) SendHelp(chatID int64) error {
	msg := tgbot.NewMessage(chatID, "I understand these commands:\n"+
		bot.Help+"\n"+
		bot.PriceByTicker+"\n"+
		bot.NewAlert)
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send help: %w", err)
	}
	return nil
}

func (b *Bot) SendTickerRequest(chatID int64) error {
	msg := tgbot.NewMessage(chatID, "Enter tickers the price of which you are interested in")
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send ticker request: %w", err)
	}
	return nil
}

func (b *Bot) SendTickerPrices(chatID int64, prices string) error {
	msg := tgbot.NewMessage(chatID, prices)
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send ticker prices: %w", err)
	}
	return nil
}

func (b *Bot) SendNameRequest(chatID int64) error {
	msg := tgbot.NewMessage(chatID, "Enter name for your alert")
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("can't send name request: %w", err)
	}
	return nil
}

func (b *Bot) SendIndicatorIdRequest(chatID int64) error {
	var indicatorKeyboard = tgbot.NewReplyKeyboard(
		tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton(domain.LowerThanValue)),
		tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton(domain.MoreThanValue)),
		tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton(domain.MoreThanMA)),
	)
	msg := tgbot.NewMessage(chatID, "Choose indicator type")
	msg.ReplyMarkup = indicatorKeyboard
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("can't indicator keyboard: %w", err)
	}
	return nil
}
