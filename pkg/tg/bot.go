package tg

import (
	"github.com/SerGGe43/alertilka/pkg/tinkoff"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

var updateConfig = tgbot.UpdateConfig{
	Offset:         0,
	Limit:          0,
	Timeout:        30,
	AllowedUpdates: nil,
}

func BotInit(token string) (tgbot.BotAPI, error) {
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return *bot, nil
}

func MainMenu(bot tgbot.BotAPI, client tinkoff.Client) {
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbot.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Text {
		case "help":
			msg.Text = "I understand \"Price by ticker\" command"
		case "Price by ticker":
			msg.Text = "Enter your tickers with spaces"
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
			updates.Clear()
			priceByTicker(bot, client)
			updates = bot.GetUpdatesChan(updateConfig)
			msg.Text = "Here's prices by your tickets"
		default:
			msg.Text = "I don't know that command"
		}
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}

func priceByTicker(bot tgbot.BotAPI, client tinkoff.Client) {
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := tgbot.NewMessage(update.Message.Chat.ID, "")
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
		msg.Text = prices_str
		if _, err := bot.Send(msg); err != nil {
			msg.Text = "bad ticker"
			_, _ = bot.Send(msg)
		}
		break
	}
	updates.Clear()
}
