package portfolio

import (
	"fmt"
	"github.com/egelis/binance/pkg/exchange"
	"github.com/shopspring/decimal"
)

// getTokenStatistic TODO: решить с постоянным GetBalance()
func (p *Portfolio) getTokenStatistic(token string) (TokenStatistic, error) {
	balance, err := p.exchange.GetBalance()
	if err != nil {
		return TokenStatistic{}, err
	}

	dividends, err := p.exchange.GetTokenDividends(token)
	if err != nil {
		return TokenStatistic{}, err
	}
	divSum := decimal.NewFromFloat(0)
	for _, div := range dividends {
		divSum = divSum.Add(div)
	}

	tradeHistory, err := p.exchange.GetPairTradeHistory(token + "USDT")
	if err != nil {
		return TokenStatistic{}, err
	}

	inStaking := getStakingAmount(tradeHistory, divSum, balance[token])

	currentPrice, err := p.exchange.GetTokenPrice(token + "USDT")
	if err != nil {
		return TokenStatistic{}, err
	}

	averagePrice := getAveragePrice(tradeHistory)

	var profit decimal.Decimal
	if averagePrice.Equal(decimal.NewFromFloat(0)) {
		profit = decimal.NewFromFloat(0)
	} else {
		profit = currentPrice.Div(averagePrice)
	}

	return TokenStatistic{
		Balance:      balance[token],
		InStaking:    inStaking,
		AveragePrice: averagePrice,
		Dividends:    divSum,
		CurrentPrice: currentPrice,
		Profit:       profit,
	}, nil
}

func getStakingAmount(
	tradeHistory []exchange.TradePoint,
	dividends decimal.Decimal,
	balance decimal.Decimal,
) decimal.Decimal {
	var amount decimal.Decimal

	for _, tradePoint := range tradeHistory {
		if tradePoint.IsBuyer {
			amount = amount.Add(tradePoint.Quantity)
			amount = amount.Sub(tradePoint.Commission)
		} else {
			amount = amount.Sub(tradePoint.Quantity)
		}
	}

	amount = amount.Add(dividends)
	amount = amount.Sub(balance)

	if amount.LessThan(decimal.NewFromFloat(0)) {
		return decimal.NewFromFloat(0)
	} else {
		return amount
	}
}

func getAveragePrice(trades []exchange.TradePoint) decimal.Decimal {
	var moneySum, quantitySum, average decimal.Decimal

	for _, point := range trades {
		if point.IsBuyer {
			moneySum = moneySum.Add(point.Quantity.Mul(point.Price))
			quantitySum = quantitySum.Add(point.Quantity)
			average = moneySum.Div(quantitySum)
		} else {
			quantitySum = quantitySum.Sub(point.Quantity)
			moneySum = average.Mul(quantitySum)
		}
	}

	return average
}

func getProfitPercents(profit decimal.Decimal) string {
	if profit.Equal(decimal.NewFromFloat(0)) {
		return fmt.Sprintf("0%%")
	}
	res := profit.Sub(decimal.NewFromFloat(1)).Mul(decimal.NewFromFloat(100))
	return fmt.Sprintf("%s%%", res.StringFixed(2))
}

func getProfitUSD(balance, inStaking, averagePrice, currentPrice decimal.Decimal) string {
	if averagePrice.Equal(decimal.NewFromFloat(0)) {
		return fmt.Sprintf("0$")
	}
	res := currentPrice.Sub(averagePrice).Mul(balance.Add(inStaking))
	return fmt.Sprintf("%v$", res.StringFixed(2))
}
