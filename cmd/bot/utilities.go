package main

import (
	"bytes"
	"fmt"
	"github.com/egelis/binance/pkg/binance"
	"github.com/shopspring/decimal"
	"log"
)

type (
	AveragePrice map[string]decimal.Decimal
)

func (ap AveragePrice) String() string {
	var buf bytes.Buffer

	for symbol, price := range ap {
		_, err := fmt.Fprintf(&buf, "Average purchase price %v: %v\n", symbol, price)
		if err != nil {
			log.Fatal(err)
		}
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

func GetAveragePrices(trades binance.TradeHistory) AveragePrice {
	prices := make(AveragePrice)

	for symbol, tradeList := range trades {
		var moneySum, quantitySum, average decimal.Decimal
		for _, point := range tradeList {
			if point.IsBuyer {
				moneySum = moneySum.Add(point.Quantity.Mul(point.Price))
				quantitySum = quantitySum.Add(point.Quantity)
				average = moneySum.Div(quantitySum)
			} else {
				quantitySum = quantitySum.Sub(point.Quantity)
				moneySum = average.Mul(quantitySum)
			}
		}

		prices[symbol] = average
	}

	return prices
}
