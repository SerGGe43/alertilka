package controller

import (
	"fmt"
	"github.com/SerGGe43/alertilka/internal/domain"
)

func (c *Controller) HandleHelp(e domain.Event) error {
	err := c.Bot.SendHelp(e.ChatId)
	if err != nil {
		return fmt.Errorf("can't send help: %w", err)
	}
	return nil
}
