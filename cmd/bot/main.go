package main

import (
	"fmt"
	"os"

	"github.com/SerGGe43/alertilka/pkg/tg"
	"github.com/SerGGe43/alertilka/pkg/tinkoff"

	"go.uber.org/zap"
)

func main() {
	fmt.Println("Hello world!")
	l, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic(err)
	}
	client, err := tinkoff.NewClient(os.Getenv("TINKOFF_TOKEN"), *l)
	if err != nil {
		panic(err)
	}
	//fmt.Println(client.GetPriceByTiker([]string{"AAPL", "TCS"}))
	bot, err := tg.BotInit(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	tg.MainMenu(bot, *client)

}
