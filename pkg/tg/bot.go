package tg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SerGGe43/alertilka/pkg/tinkoff"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var updateConfig = tgbot.UpdateConfig{
	Offset:         0,
	Limit:          0,
	Timeout:        30,
	AllowedUpdates: nil,
}

var Bot tgbot.BotAPI

var mainMenuKeyboard = tgbot.NewReplyKeyboard(
	tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("help")),
	tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("Price by ticker")),
)

var updates tgbot.UpdatesChannel

func BotInit(token string) (tgbot.BotAPI, error) {
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	updates = bot.GetUpdatesChan(updateConfig)
	return *bot, nil
}

func MainMenu(bot tgbot.BotAPI, client tinkoff.Client) {
	Bot = bot
	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		chooseKeyboard("MainMenu", chatID)
		commandHandler(update.Message.Text, client, chatID)
	}
}

func priceByTicker(client tinkoff.Client, chatID int64) {
	_ = sendMessage("Enter your tickers with spaces", chatID)
	for update := range updates {
		if update.Message == nil || update.Message.Chat.ID != chatID {
			continue
		}
		tickers := strings.Split(update.Message.Text, " ")
		prices, err := client.GetPriceByTiker(tickers)
		if err != nil {
			panic(err)
		}
		prices_str := ""
		for i := range prices {
			prices_str += strconv.Itoa(prices[i])
			prices_str += " "
		}
		_ = sendMessage(prices_str, update.Message.Chat.ID)
		_ = sendMessage("Here's prices by your tickets", chatID)
		break
	}
	updates.Clear()
}

func commandHandler(text string, client tinkoff.Client, chatID int64) {
	switch text {
	case "help":
		help(chatID)
	case "Price by ticker":
		priceByTicker(client, chatID)
	default:
		_ = sendMessage("I don't know that command", chatID)
	}
}

func sendMessage(text string, chatID int64) error {
	msg := tgbot.NewMessage(chatID, text)
	if _, err := Bot.Send(msg); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func chooseKeyboard(keyboardType string, chatID int64) {
	switch keyboardType {
	case "MainMenu":
		_ = sendMainMenuKeyboard(chatID)
	default:
	}
}

func help(chatID int64) {
	//Тут в теории потом можно будет передавать список команд
	_ = sendMessage("I understand \"Price by ticker\" command", chatID)
}

func sendMainMenuKeyboard(chatID int64) error {
	msg := tgbot.NewMessage(chatID, "")
	msg.ReplyMarkup = mainMenuKeyboard
	if _, err := Bot.Send(msg); err != nil {
		return fmt.Errorf("can't send keyboard: %w", err)
	}
	return nil
}
