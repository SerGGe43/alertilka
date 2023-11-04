package controller

import (
	"fmt"
	"github.com/SerGGe43/alertilka/internal/domain"
	"strings"
)

func (c *Controller) HandleNewAlert(e domain.Event) error {
	err := c.Bot.SendTickerRequest(e.ChatId)
	if err != nil {
		return fmt.Errorf("can't send ticker request: %w", err)
	}
	err = c.UserDB.SetState(e.ChatId, domain.NEW_ALERT)
	if err != nil {
		return fmt.Errorf("can't set state new alert: %w", err)
	}
	return nil
}

func (c *Controller) newAlertHandler(e domain.Event) error {
	// TODO придумать, как ловить следующие сообщения, чтобы добавить имя, индикаторid
	ticker := strings.ToUpper(e.Message)
	user, err := c.UserDB.GetByChatId(e.ChatId)
	if err != nil {
		return fmt.Errorf("can't get user in newAlertHandler: %w", err)
	}
	alertID, err := c.AlertDB.Add(domain.Alert{
		Ticker: ticker,
		Name:   "",
		UserID: user.Id,
	})
	if err != nil {
		return fmt.Errorf("can't add alert: %w", err)
	}
	err = c.addIndicator(alertID)
	if err != nil {
		return fmt.Errorf("can't add indicator: %w", err)
	}
	err = c.Bot.SendNameRequest(e.ChatId)
	if err != nil {
		return fmt.Errorf("can't send name request: %w", err)
	}
	err = c.UserDB.SetState(e.ChatId, domain.ADD_NAME_TO_ALERT)
	if err != nil {
		return fmt.Errorf("can't set state add name to alert: %w", err)
	}
	return nil
}

func (c *Controller) addAlertName(e domain.Event) error {
	user, err := c.UserDB.GetByChatId(e.ChatId)
	if err != nil {
		return fmt.Errorf("can't get user by chatID in addAlertName: %w", err)
	}
	err = c.AlertDB.AddName(e.Message, user.Id)
	fmt.Println(c.AlertDB.GetByUserID(user.Id))
	if err != nil {
		return fmt.Errorf("can't set name: %w", err)
	}
	err = c.UserDB.SetState(e.ChatId, domain.ADD_INDICATOR_ID)
	if err != nil {
		return fmt.Errorf("can't set state add indicator id: %w", err)
	}
	return nil
}
