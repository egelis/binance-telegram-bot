package binance

import (
	"context"
	"github.com/egelis/binance/pkg/exchange"
	"github.com/shopspring/decimal"
)

func (c *Client) GetBalance() (exchange.Balance, error) {
	client, err := c.binanceAPIClient.NewGetAccountService().
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	balance := make(exchange.Balance)
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

func (c *Client) GetPairTradeHistory(pair string) ([]exchange.TradePoint, error) {
	tradesList, err := c.binanceAPIClient.NewListTradesService().
		Symbol(pair).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var res []exchange.TradePoint
	for _, tpv3 := range tradesList {
		tradePoint, err := getTradePoint(tpv3)
		if err != nil {
			return nil, err
		}

		res = append(res, tradePoint)
	}

	return res, nil
}

// GetTokenDividends TODO: убрать зависимость от Limit(500)
func (c *Client) GetTokenDividends(token string) (exchange.Dividends, error) {
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
