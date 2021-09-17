package binance

import (
	"context"
	"fmt"
	binanceApi "github.com/adshao/go-binance/v2"
	"strconv"
)

type trade struct {
	price      float64
	quantity   float64
	commission float64
	isBuyer    bool
}

type Client struct {
	client *binanceApi.Client
}

func NewClient(apiKey, secretKey string) *Client {
	client := binanceApi.NewClient(apiKey, secretKey)
	client.NewSetServerTimeService().Do(context.Background())

	return &Client{client: client}
}

func (c *Client) PrintAveragePrices(trades map[string][]trade) {
	for symbol, tradeList := range trades {
		var moneySum, quantitySum, average float64
		for _, point := range tradeList {
			if point.isBuyer {
				moneySum += point.quantity * point.price
				quantitySum += point.quantity
				average = moneySum / quantitySum
			} else {
				quantitySum -= point.quantity
				moneySum = average * quantitySum
			}
		}

		fmt.Printf("Average purchase price %v: %.8f\n", symbol, average)
	}
}

func (c *Client) GetRubCourse() (float64, error) {
	priceStr, err := c.client.NewListPricesService().Symbol("USDTRUB").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	price, err := strconv.ParseFloat(priceStr[0].Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func (c *Client) GetTradeHistory(pairs []string) (map[string][]trade, error) {
	trades := make(map[string][]trade)

	for _, pair := range pairs {
		tradesList, err := c.client.NewListTradesService().Symbol(pair).
			Do(context.Background())

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for _, tradePoint := range tradesList {
			price, _ := strconv.ParseFloat(tradePoint.Price, 64)
			quantity, _ := strconv.ParseFloat(tradePoint.Quantity, 64)
			commission, _ := strconv.ParseFloat(tradePoint.Commission, 64)
			isBuyer := tradePoint.IsBuyer

			trades[tradePoint.Symbol] = append(trades[tradePoint.Symbol],
				trade{
					price:      price,
					quantity:   quantity,
					commission: commission,
					isBuyer:    isBuyer,
				})
		}
	}

	return trades, nil
}

func (c *Client) MakeTokenPairsFromBalance(b map[string]float64) []string {
	var pairs []string
	for coin := range b {
		if coin != "USDT" {
			pairs = append(pairs, coin+"USDT")
		}
	}
	return pairs
}

func (c *Client) GetBalance() (map[string]float64, error) {
	client, err := c.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	balance := make(map[string]float64)
	for _, coin := range client.Balances {
		free, err := strconv.ParseFloat(coin.Free, 64)

		if err != nil {
			return nil, err
		}

		if free != 0 {
			balance[coin.Asset] = free
		}
	}

	return balance, nil
}
