package controller

import (
	"fmt"
	"github.com/SerGGe43/alertilka/internal/domain"
	"strconv"
	"strings"
)

func (c *Controller) HandlePrice(e domain.Event) error {
	err := c.Bot.SendTickerRequest(e.ChatId)
	if err != nil {
		return fmt.Errorf("can't send ticker request: %w", err)
	}
	err = c.UserDB.SetState(e.ChatId, domain.PRICE_BY_TICKER)
	if err != nil {
		return fmt.Errorf("can't set state on price by ticker: %w", err)
	}
	return nil
}

func (c *Controller) tickerHandler(e domain.Event) error {
	tickers := strings.Split(e.Message, " ")
	prices, err := c.Client.GetPriceByTiker(tickers)
	if err != nil {
		panic(err)
	}
	prices_str := ""
	for i := range prices {
		prices_str += strconv.Itoa(prices[i])
		prices_str += "$ "
	}
	err = c.Bot.SendTickerPrices(e.ChatId, prices_str)
	if err != nil {
		return fmt.Errorf("Can't send ticker prices: %w", err)
	}
	err = c.UserDB.SetState(e.ChatId, domain.MAIN_MENU)
	if err != nil {
		return fmt.Errorf("Can't set main menu before price by ticker: %w", err)
	}
	return nil
}
