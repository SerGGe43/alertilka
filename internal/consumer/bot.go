package consumer

import (
	"fmt"
	"github.com/SerGGe43/alertilka/internal/domain"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Consumer struct {
	api *tgbot.BotAPI
}

func NewConsumer(api *tgbot.BotAPI) (*Consumer, error) {
	return &Consumer{api: api}, nil
}

func (c *Consumer) Consume() (<-chan domain.Event, <-chan error) {
	tgChan := c.api.GetUpdatesChan(tgbot.UpdateConfig{
		Offset:         0,
		Limit:          0,
		Timeout:        30,
		AllowedUpdates: nil,
	})

	resChan, errChan := make(chan domain.Event), make(chan error)

	go func() {
		for {
			select {
			case update, ok := <-tgChan:
				if !ok {
					errChan <- fmt.Errorf("chanel was closed")
				}

				resChan <- updateToDomain(update)
			}
		}
	}()

	return resChan, errChan
}

func updateToDomain(update tgbot.Update) domain.Event {
	return domain.Event{
		// TODO: check nil
		ChatId:  update.Message.Chat.ID,
		Message: update.Message.Text,
		Name:    update.Message.Chat.UserName,
	}
}
