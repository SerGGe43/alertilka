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
	ticker := strings.ToUpper(e.Message)
	user, err := c.UserDB.GetByChatId(e.ChatId)
	if err != nil {
		return fmt.Errorf("can't get user in newAlertHandler: %w", err)
	}
	_, err = c.AlertDB.Add(domain.Alert{
		Ticker: ticker,
		Name:   "",
		UserID: user.Id,
	})
	if err != nil {
		return fmt.Errorf("can't add alert: %w", err)
	}
	return nil
}

func (c *Controller) addAlertName(e domain.Event) error {
	return nil
}
