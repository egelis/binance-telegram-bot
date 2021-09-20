package main

import (
	"fmt"
	"github.com/egelis/binance/pkg/binance"
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
			//amount -= tradePoint.Quantity
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

func printTokensStatistic(c *binance.Client) {
	balance, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}

	pairs := GetTokenPairs(TOKENS)
	tradeHistory, err := c.GetTradeHistory(pairs)
	if err != nil {
		log.Fatal(err)
	}

	averagePrice := GetAveragePrices(tradeHistory)

	for _, token := range TOKENS {
		dividends, err := c.GetTokenDividends(token)
		if err != nil {
			log.Fatal(err)
		}

		divSum := decimal.NewFromFloat(0)
		for _, div := range dividends {
			divSum = divSum.Add(div)
		}

		staking := getStakingAmount(tradeHistory[token+"USDT"], divSum, balance[token])

		fmt.Printf("Token: %s\n"+
			"Balance: %v\n"+
			"In staking: %v\n"+
			"Average: %v\n"+
			"Dividends: %v\n\n",
			token, balance[token], staking, averagePrice[token+"USDT"], divSum)
	}
}
