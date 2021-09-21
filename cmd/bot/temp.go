package main

import (
	"fmt"
	"github.com/egelis/binance/pkg/exchange/binance"
	"github.com/shopspring/decimal"
	"log"
)

// Token: KSM
// Должно быть: 0.33033445
// Имеется: 0.33030195
// Разница: 0.0000325 откуда??

func getStakingAmount(tradeHistory []binance.TradePoint, divSum decimal.Decimal, balance decimal.Decimal) decimal.Decimal {
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

func GetTokensStatistic(c *binance.Client) {
	balance, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}

	for _, token := range TOKENS {
		dividends, err := c.GetTokenDividends(token)
		if err != nil {
			log.Fatal(err)
		}

		divSum := decimal.NewFromFloat(0)
		for _, div := range dividends {
			divSum = divSum.Add(div)
		}

		tradeHistory, err := c.GetPairTradeHistory(token + "USDT")
		if err != nil {
			log.Fatal(err)
		}

		staking := getStakingAmount(tradeHistory, divSum, balance[token])

		currentPrice, err := c.GetTokenPrice(token + "USDT")
		if err != nil {
			log.Fatal(err)
		}

		averagePrice := GetAveragePrice(tradeHistory)
		var profit decimal.Decimal
		if !averagePrice.Equal(decimal.NewFromFloat(0)) {
			profit = currentPrice.Div(averagePrice)
		} else {
			profit = decimal.NewFromFloat(0)
		}

		fmt.Printf("Token:         %s\n"+
			"Balance:       %v\n"+
			"In staking:    %v\n"+
			"Average:       %v\n"+
			"Dividends:     %v\n"+
			"Current price: %v\n"+
			"Profit:        %v\n\n",
			token,
			balance[token].StringFixed(8),
			staking.StringFixed(8),
			averagePrice.StringFixed(2),
			divSum.StringFixed(8),
			currentPrice.StringFixed(2),
			profit.StringFixed(4))
	}
}

/*
return TokenStatistic{
Token:        token,
OnBalance:    balance[token],
InStaking:    staking,
Average:      averagePrice,
Dividends:    divSum,
CurrentPrice: currentPrice,
Profit:       profit,
}
*/
