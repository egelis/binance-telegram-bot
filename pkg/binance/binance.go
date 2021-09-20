package binance

import (
	"bytes"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"log"

	binanceAPI "github.com/adshao/go-binance/v2"
)

type (
	Balance      map[string]decimal.Decimal
	TradeHistory map[string][]TradePoint
	Dividend     []decimal.Decimal

	TradePoint struct {
		Price           decimal.Decimal
		Quantity        decimal.Decimal
		Commission      decimal.Decimal
		CommissionAsset string
		IsBuyer         bool
	}

	Client struct {
		binanceAPIClient *binanceAPI.Client
	}
)

func NewClient(apiKey, secretKey string) (*Client, error) {
	newClient := binanceAPI.NewClient(apiKey, secretKey)

	if _, err := newClient.NewSetServerTimeService().Do(context.Background()); err != nil {
		return nil, nil
	}

	return &Client{binanceAPIClient: newClient}, nil
}

func (b Balance) String() string {
	var buf bytes.Buffer

	_, err := fmt.Fprintln(&buf, "Balance:")
	if err != nil {
		log.Fatal(err)
	}

	for token, value := range b {
		_, err = fmt.Fprintf(&buf, "%s: %f\n", token, value)
		if err != nil {
			log.Fatal(err)
		}
	}

	return buf.String()
}
