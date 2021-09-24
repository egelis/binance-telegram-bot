package portfolio

import (
	"github.com/egelis/binance/pkg/exchange"
	"github.com/shopspring/decimal"
)

func getStakingAmount(tradeHistory []exchange.TradePoint, divSum decimal.Decimal, balance decimal.Decimal) decimal.Decimal {
	var amount decimal.Decimal

	for _, tradePoint := range tradeHistory {
		if tradePoint.IsBuyer {
			amount = amount.Add(tradePoint.Quantity)
			amount = amount.Sub(tradePoint.Commission)
		} else {
			amount = amount.Sub(tradePoint.Quantity)
		}
	}

	amount = amount.Add(divSum)
	amount = amount.Sub(balance)

	if amount.LessThan(decimal.NewFromFloat(0)) {
		return decimal.NewFromFloat(0)
	} else {
		return amount
	}
}

func getAveragePrices(trades []exchange.TradePoint) decimal.Decimal {
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
