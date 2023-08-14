package tinkoff

import (
	"context"
	"fmt"
	sdk "github.com/tinkoff/invest-api-go-sdk/investgo"
	investapi "github.com/tinkoff/invest-api-go-sdk/proto"
	"go.uber.org/zap"
)

type Client struct {
	client      *sdk.Client
	tikerToFigi map[string]string
}

func NewClient(token string, logger zap.Logger) (*Client, error) {
	client, err := sdk.NewClient(
		context.TODO(),
		sdk.Config{
			EndPoint:                      "sandbox-invest-public-api.tinkoff.ru:443",
			Token:                         token,
			AppName:                       "Alertilka",
			AccountId:                     "1",
			DisableResourceExhaustedRetry: false,
			DisableAllRetry:               false,
			MaxRetries:                    1,
		},
		logger.Sugar(),
	)
	if err != nil {
		return nil, fmt.Errorf("can't create client: %w", err)
	}
	assets, err := client.NewInstrumentsServiceClient().GetAssets()
	if err != nil {
		return nil, fmt.Errorf("can't get assets: %w", err)
	}
	tikers := generateTickersMap(assets.Assets)
	return &Client{
		client:      client,
		tikerToFigi: tikers,
	}, nil
}

func (c Client) GetPriceByTiker(tickers []string) ([]int, error) {
	figi := make([]string, 0)
	for _, t := range tickers {
		figi = append(figi, c.tikerToFigi[t])
	}
	lastPrices, err := c.client.NewMarketDataServiceClient().GetLastPrices(figi)
	if err != nil {
		return nil, fmt.Errorf("can't get last prices: %w", err)
	}
	prices := make([]int, 0)
	for _, p := range lastPrices.LastPrices {
		prices = append(prices, int(p.Price.Units))
	}
	return prices, nil
}

func generateTickersMap(assets []*investapi.Asset) map[string]string {
	tikers := make(map[string]string)
	for _, e := range assets {
		for _, i := range e.Instruments {
			tikers[i.Ticker] = i.Figi
		}
	}
	return tikers
}
