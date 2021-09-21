package binance

import (
	"context"
	"strconv"

	"github.com/shopspring/decimal"
)

func (c *Client) GetTradeHistoryForPair(pair string) ([]TradePoint, error) {
	tradesList, err := c.binanceAPIClient.NewListTradesService().Symbol(pair).Do(context.Background())
	if err != nil {
		return nil, err
	}

	var res []TradePoint
	for _, tpv3 := range tradesList {
		point, err := getTradePoint(tpv3)
		if err != nil {
			return nil, err
		}

		res = append(res, point)
	}

	return res, nil
}

func (c *Client) GetTradeHistoryForPairs(pairs []string) (TradeHistory, error) {
	trades := make(TradeHistory)

	for _, pair := range pairs {
		tradePoints, _ := c.GetTradeHistoryForPair(pair)
		trades[pair] = tradePoints
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
		free, err := decimal.NewFromString(coin.Free)
		if err != nil {
			return nil, err
		}

		if !free.Equal(decimal.NewFromFloat(0)) {
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

func (c *Client) GetTokenDividends(token string) ([]decimal.Decimal, error) {
	dividends, err := c.binanceAPIClient.NewAssetDividendService().
		Limit(500).
		Asset(token).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	res := make([]decimal.Decimal, 0, len(*dividends.Rows))
	for _, row := range *dividends.Rows {
		amount, err := decimal.NewFromString(row.Amount)
		if err != nil {
			return nil, err
		}

		res = append(res, amount)
	}

	return res, nil
}

func (c *Client) GetTokenPrice(symbol string) (decimal.Decimal, error) {
	tokenPrice, err := c.binanceAPIClient.NewListPricesService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return decimal.Decimal{}, err
	}

	res, err := decimal.NewFromString(tokenPrice[0].Price)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return res, err
}
