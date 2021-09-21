package binance

import (
	"context"
	binanceAPI "github.com/adshao/go-binance/v2"
)

type Client struct {
	binanceAPIClient *binanceAPI.Client
}

func NewClient(apiKey, secretKey string) (*Client, error) {
	newClient := binanceAPI.NewClient(apiKey, secretKey)

	if _, err := newClient.NewSetServerTimeService().Do(context.Background()); err != nil {
		return nil, nil
	}

	return &Client{binanceAPIClient: newClient}, nil
}
