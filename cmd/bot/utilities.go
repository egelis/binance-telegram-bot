package main

import (
	"bytes"
	"fmt"
	"github.com/egelis/binance/pkg/exchange/binance"
	"github.com/shopspring/decimal"
	"log"
)

type (
	AveragePrices map[string]decimal.Decimal

	TokenStatistic struct {
		Token        string
		OnBalance    decimal.Decimal
		InStaking    decimal.Decimal
		Average      decimal.Decimal
		Dividends    decimal.Decimal
		CurrentPrice decimal.Decimal
		Profit       decimal.Decimal
	}
)

func (ap AveragePrices) String() string {
	var buf bytes.Buffer

	for symbol, price := range ap {
		_, err := fmt.Fprintf(&buf, "Average purchase price %s: %v\n", symbol, price.StringFixed(2))
		if err != nil {
			log.Fatal(err)
		}
	}

	return buf.String()
}

func (ts TokenStatistic) String() string {
	var buf bytes.Buffer

	_, err := fmt.Fprintf(
		&buf, "Token:         %s\n"+
			"Balance:       %v\n"+
			"In staking:    %v\n"+
			"Average:       %v\n"+
			"Dividends:     %v\n"+
			"Current price: %v\n"+
			"Profit:        %v\n",
		ts.Token,
		ts.OnBalance.StringFixed(8),
		ts.InStaking.StringFixed(8),
		ts.Average.StringFixed(2),
		ts.Dividends.StringFixed(8),
		ts.CurrentPrice.StringFixed(2),
		ts.Profit.StringFixed(4),
	)

	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func GetTokenPairs(tokens []string) []string {
	var pairs []string

	for _, coin := range tokens {
		if coin != "USDT" {
			pairs = append(pairs, coin+"USDT")
		}
	}

	return pairs
}

func GetAveragePrices(tradeHistory binance.TradeHistory) AveragePrices {
	prices := make(AveragePrices)

	for symbol, tokensTradeList := range tradeHistory {
		prices[symbol] = GetAveragePrice(tokensTradeList)
	}

	return prices
}

func GetAveragePrice(tokensTradeList []binance.TradePoint) decimal.Decimal {
	var moneySum, quantitySum, average decimal.Decimal
	for _, tradePoint := range tokensTradeList {
		if tradePoint.IsBuyer {
			moneySum = moneySum.Add(tradePoint.Quantity.Mul(tradePoint.Price))
			quantitySum = quantitySum.Add(tradePoint.Quantity)
			average = moneySum.Div(quantitySum)
		} else {
			quantitySum = quantitySum.Sub(tradePoint.Quantity)
			moneySum = average.Mul(quantitySum)
		}
	}

	return average
}
