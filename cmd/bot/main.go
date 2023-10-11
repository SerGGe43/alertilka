package main

import (
	"context"
	"fmt"
	"github.com/SerGGe43/alertilka/internal/consumer"
	"github.com/SerGGe43/alertilka/internal/controller"
	"github.com/SerGGe43/alertilka/internal/postgres"
	"github.com/SerGGe43/alertilka/internal/tg"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pressly/goose/v3"
	"os"

	"github.com/SerGGe43/alertilka/pkg/tinkoff"

	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {
	var host = os.Getenv("HOST")
	var port = os.Getenv("PORT")
	var user = os.Getenv("USER")
	var password = os.Getenv("PASSWORD")
	var dbname = os.Getenv("DBNAME")
	var sslmode = os.Getenv("SSLMODE")
	var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	fmt.Println("Hello world!")
	l, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		return fmt.Errorf("can't create logger for tinkoff client: %w", err)
	}
	client, err := tinkoff.NewClient(os.Getenv("TINKOFF_TOKEN"), *l)
	if err != nil {
		return fmt.Errorf("can't create tinkoff client: %w", err)
	}

	db, err := postgres.NewConnection(dbInfo)
	if err != nil {
		return fmt.Errorf("can't create postgres connection: %w", err)
	}
	userDB := postgres.NewUser(db)
	err = goose.Up(db, "migrations")

	api, err := tgbot.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		return fmt.Errorf("can't create bot api: %w", err)
	}
	bot := tg.NewBot(api)

	cons, err := consumer.NewConsumer(api)
	if err != nil {
		return fmt.Errorf("can't create consumer: %w", err)
	}
	eventChan, _ := cons.Consume()

	ctrl := controller.NewController(userDB, bot, *client)
	err = ctrl.Run(context.TODO(), eventChan)
	if err != nil {
		return fmt.Errorf("can't run controller: %w", err)
	}
	return nil
}
