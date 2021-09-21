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

func printTokensStatistic(c *binance.Client) {
	balance, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}

	pairs := GetTokenPairs(TOKENS)
	tradeHistory, err := c.GetTradeHistoryForPairs(pairs)
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

		tokenUSDT := token + "USDT"

		staking := getStakingAmount(tradeHistory[tokenUSDT], divSum, balance[token])

		currentPrice, err := c.GetTokenPrice(tokenUSDT)
		if err != nil {
			log.Fatal(err)
		}

		var profit decimal.Decimal
		if !averagePrice[tokenUSDT].Equal(decimal.NewFromFloat(0)) {
			profit = currentPrice.Div(averagePrice[tokenUSDT])
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
			averagePrice[tokenUSDT].StringFixed(2),
			divSum.StringFixed(8),
			currentPrice.StringFixed(2),
			profit.StringFixed(4))
	}
}
