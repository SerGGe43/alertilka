package controller

import (
	"fmt"
	"github.com/SerGGe43/alertilka/internal/domain"
)

func (c *Controller) HandleStart(e domain.Event) error {
	user, err := c.UserDB.GetByChatId(e.ChatId)
	if err != nil {
		return fmt.Errorf("can't get user: %w", err)
	}
	if user == nil {
		_, err = c.UserDB.Add(domain.User{
			Name:   e.Name,
			ChatId: e.ChatId,
		})
		if err != nil {
			return fmt.Errorf("can't add user to db: %w", err)
		}
	}
	err = c.Bot.SendMenu(e.ChatId)
	if err != nil {
		return fmt.Errorf("can't send menu: %w", err)
	}
	return nil
}
