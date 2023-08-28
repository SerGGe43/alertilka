package tg

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/SerGGe43/alertilka/migrations"
	"github.com/SerGGe43/alertilka/pkg/tinkoff"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
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
	tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("New alert")),
)

var trackingTypeKeyboard = tgbot.NewReplyKeyboard(
	tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("more than price")),
	tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("lower than price")),
	tgbot.NewKeyboardButtonRow(tgbot.NewKeyboardButton("more than MA")),
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

var updates tgbot.UpdatesChannel

func BotInit(token string) (tgbot.BotAPI, error) {
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	updates = bot.GetUpdatesChan(updateConfig)
	return *bot, nil
}

func migration() {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("can't open sql: %w", err)
	}
	err = goose.Up(db, ".")
	if err != nil {
		fmt.Printf("can't migrate: %w", err)
	}
}

func addUser(chatID int64, username string) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("can't open sql: %w", err)
	}
	//INSERT INTO "User"(name, chatid, tgid) VALUES($1, $2, $3);
	defer db.Close()
	if chats := db.QueryRow(`SELECT chatid FROM "User" WHERE chatid = $1`, chatID).Scan(); chats != nil {
		if chats == sql.ErrNoRows {
			data := `INSERT INTO "User"(name, chatid) VALUES($1, $2);`
			if _, err := db.Exec(data, `@`+username, chatID); err != nil {
				fmt.Printf("can't insert: %w", err)
			}
		}
		return
	}
}

func addTickerToTrack(chatID int64, ticker, tickerName string) {
	var alertID int
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("can't open sql: %w", err)
	}
	defer db.Close()
	data := `INSERT INTO alert(ticker, name, userID) VALUES($1, $2, (SELECT id FROM "User" WHERE chatID = $3));`
	if _, err := db.Exec(data, ticker, tickerName, chatID); err != nil {
		fmt.Printf("can't insert ticker to track: %w", err)
	}
	//alertID, err := db.Exec(`SELECT id FROM alert WHERE userID = $1 AND name = $2, AND ticker = $3`, chatID, tickerName, ticker)
	db.QueryRow(`SELECT id FROM alert WHERE userID = $1 AND name = $2, AND ticker = $3`, chatID, tickerName, ticker).Scan(&alertID)
	fmt.Println(alertID)
	selectTrackingType(chatID, tickerName)
	_ = sendMessage("Your alert successfully added", chatID)
}

func addTrackingType(trackingType, name string, value int, chatID int64) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Printf("can't open sql: %w", err)
	}
	defer db.Close()
	data := `INSERT INTO indicator(alertid, indicatorID, value) VALUES((SELECT id FROM alert 
                                                                     WHERE userID = (SELECT id FROM "User" WHERE chatID = $1) AND name = $2), $3, $4);`
	if _, err := db.Exec(data, chatID, name, trackingType, value); err != nil {
		fmt.Printf("can't insert: %w", err)
	}
}

func MainMenu(bot tgbot.BotAPI, client tinkoff.Client) {
	Bot = bot
	migration()
	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		username := update.Message.Chat.UserName
		//fmt.Println(username, tgid)
		addUser(chatID, username)
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
	case "/start":
		_ = sendMessage("Hello, I'm Alertilka bot.\n "+
			"You can use \"help\" command to find out what i can do.", chatID)
		chooseKeyboard("MainMenu", chatID)
	case "help":
		help(chatID)
	case "Price by ticker":
		priceByTicker(client, chatID)
	case "New alert":
		newAlert(chatID)
	default:
		_ = sendMessage("I don't know that command", chatID)
		chooseKeyboard("MainMenu", chatID)
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
	case "TrackingType":
		_ = sendTrackingTypeKeyboard(chatID)
	default:
	}
}

func help(chatID int64) {
	//Тут в теории потом можно будет передавать список команд
	_ = sendMessage("I understand \"Price by ticker\" command", chatID)
}

func sendMainMenuKeyboard(chatID int64) error {
	msg := tgbot.NewMessage(chatID, "menu keyboard")
	//msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
	//if _, err := Bot.Send(msg); err != nil {
	//	return fmt.Errorf("can't send keyboard: %w", err)
	//}
	msg.ReplyMarkup = mainMenuKeyboard
	if _, err := Bot.Send(msg); err != nil {
		return fmt.Errorf("can't send keyboard: %w", err)
	}
	return nil
}
func sendTrackingTypeKeyboard(chatID int64) error {
	msg := tgbot.NewMessage(chatID, "Choose tracking type between \"more than price\", \"lower than price\", \"more than MA\"")
	//msg.ReplyMarkup = tgbot.NewRemoveKeyboard(true)
	//if _, err := Bot.Send(msg); err != nil {
	//	return fmt.Errorf("can't send keyboard: %w", err)
	//}
	msg.ReplyMarkup = trackingTypeKeyboard
	if _, err := Bot.Send(msg); err != nil {
		return fmt.Errorf("can't send keyboard: %w", err)
	}
	return nil
}

func newAlert(chatID int64) {
	var ticker, tickerName string
	_ = sendMessage("Enter the ticker you want to track (4 symbols)", chatID)
	for update := range updates {
		if update.Message == nil || update.Message.Chat.ID != chatID {
			continue
		}
		ticker = update.Message.Text
		strings.ToUpper(ticker)
		break
	}
	_ = sendMessage("Enter the name for your ticker", chatID)
	for update := range updates {
		if update.Message == nil || update.Message.Chat.ID != chatID {
			continue
		}
		tickerName = update.Message.Text
		break
	}
	addTickerToTrack(chatID, ticker, tickerName)
}

func selectTrackingType(chatID int64, name string) {
	var trackingType string
	var value int
	var err error
	chooseKeyboard("TrackingType", chatID)
	for update := range updates {
		if update.Message == nil || update.Message.Chat.ID != chatID {
			continue
		}
		trackingType = update.Message.Text
		//strings.ToLower(trackingType)
		fmt.Println(trackingType)
		break
	}
	_ = sendMessage("Enter value to compare", chatID)
	for update := range updates {
		if update.Message == nil || update.Message.Chat.ID != chatID {
			continue
		}
		value, err = strconv.Atoi(update.Message.Text)
		if err != nil {
			fmt.Printf("can't convert value %w", err)
		}
		break
	}
	addTrackingType(trackingType, name, value, chatID)
}
