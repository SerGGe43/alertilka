package controller

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/SerGGe43/alertilka/internal/bot"
	"github.com/SerGGe43/alertilka/internal/domain"
	"github.com/SerGGe43/alertilka/internal/repository"
	"github.com/SerGGe43/alertilka/pkg/tinkoff"
	"log"
)

type Controller struct {
	UserDB      repository.User
	AlertDB     repository.Alert
	IndicatorDB repository.Indicator
	Bot         bot.Bot
	Client      tinkoff.Client
}

func NewController(userDB repository.User, alertDB repository.Alert, indicatorDB repository.Indicator,
	bot bot.Bot, client tinkoff.Client) *Controller {
	return &Controller{
		UserDB:      userDB,
		AlertDB:     alertDB,
		IndicatorDB: indicatorDB,
		Bot:         bot,
		Client:      client,
	}
}

func (c *Controller) Run(ctx context.Context, event <-chan domain.Event) error {
	for {
		select {
		case e := <-event:
			err := c.stateGetter(e)
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			return fmt.Errorf("context closed")
		}
	}
}

func (c *Controller) stateGetter(e domain.Event) error {
	state, err := c.UserDB.GetState(e.ChatId)
	if err != nil {
		if err == sql.ErrNoRows {
			state = 0
			err = c.stateHandler(state, e)
			if err != nil {
				return fmt.Errorf("can't handle state: %w", err)
			}
			return nil
		}
		return fmt.Errorf("can't get state: %w", err)
	}
	err = c.stateHandler(state, e)
	if err != nil {
		return fmt.Errorf("can't handle state: %w", err)
	}
	return nil
}

func (c *Controller) stateHandler(state int64, e domain.Event) error {
	var err error
	switch state {
	case domain.MAIN_MENU:
		err = c.mainMenuHandler(e)
	case domain.PRICE_BY_TICKER:
		err = c.tickerHandler(e)
	case domain.NEW_ALERT:
		err = c.newAlertHandler(e)
	case domain.ADD_NAME_TO_ALERT:
		err = c.addAlertName(e)
	case domain.ADD_INDICATOR_ID:
		err = c.addIndicatorId(e)
	}

	//default:
	//	_ = sendMessage("I don't know that command", chatID)
	//	chooseKeyboard("MainMenu", chatID)
	if err != nil {
		return fmt.Errorf("can't handle event: %w", err)
	}
	return nil
}

func (c *Controller) mainMenuHandler(e domain.Event) error {
	var err error
	switch e.Message {
	case bot.Start:
		err = c.HandleStart(e)
	case bot.Help:
		err = c.HandleHelp(e)
	case bot.PriceByTicker:
		err = c.HandlePrice(e)
	case bot.NewAlert:
		err = c.HandleNewAlert(e)
	}
	if err != nil {
		return fmt.Errorf("can't handle main menu: %w", err)
	}
	return nil
}
