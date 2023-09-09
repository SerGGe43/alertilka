package controller

import (
	"context"
	"fmt"
	"github.com/SerGGe43/alertilka/internal/bot"
	"github.com/SerGGe43/alertilka/internal/domain"
	"github.com/SerGGe43/alertilka/internal/repository"
	"log"
)

type Controller struct {
	UserDB repository.User
	Bot    bot.Bot
}

func NewController(userDB repository.User, bot bot.Bot) *Controller {
	return &Controller{
		UserDB: userDB,
		Bot:    bot,
	}
}

func (c *Controller) Run(ctx context.Context, event <-chan domain.Event) error {
	for {
		select {
		case e := <-event:
			err := c.commandHandler(e)
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			return fmt.Errorf("context closed")
		}
	}
}

func (c *Controller) commandHandler(e domain.Event) error {
	var err error
	switch e.Message {
	case bot.Start:
		err = c.HandleStart(e)
		//case "help":
		//	help(chatID)
		//case "Price by ticker":
		//	priceByTicker(client, chatID)
		//case "New alert":
		//	newAlert(chatID, dbInfo)
		//default:
		//	_ = sendMessage("I don't know that command", chatID)
		//	chooseKeyboard("MainMenu", chatID)
	}
	if err != nil {
		return fmt.Errorf("can't handle event: %w", err)
	}
	return nil
}
