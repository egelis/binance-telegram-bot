package binance

import (
	"context"
	"strconv"
)

func (c *Client) GetTradeHistory(pairs []string) (TradeHistory, error) {
	trades := make(TradeHistory)

	for _, pair := range pairs {
		tradesList, err := c.binanceAPIClient.NewListTradesService().Symbol(pair).Do(context.Background())
		if err != nil {
			return nil, err
		}

		for _, tpv3 := range tradesList {
			point, err := getTradePoint(tpv3)
			if err != nil {
				return nil, err
			}

			trades[tpv3.Symbol] = append(trades[tpv3.Symbol], point)
		}
	}

	return trades, nil
}

func (c *Client) GetBalance() (Balance, error) {
	client, err := c.binanceAPIClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	balance := make(Balance)
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

func (c *Client) GetRubCourse() (float64, error) {
	priceStr, err := c.binanceAPIClient.NewListPricesService().Symbol("USDTRUB").Do(context.Background())
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(priceStr[0].Price, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func (c *Client) GetTokenDividends(token string) ([]float64, error) {
	dividends, err := c.binanceAPIClient.NewAssetDividendService().
		Limit(500).
		Asset(token).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	res := make([]float64, 0, len(*dividends.Rows))
	for _, row := range *dividends.Rows {
		amount, err := strconv.ParseFloat(row.Amount, 64)
		if err != nil {
			return nil, err
		}

		res = append(res, amount)
	}

	return res, nil
}
