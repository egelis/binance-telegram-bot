package main

import (
	"bytes"
	"fmt"
	"github.com/egoreli/binance/pkg/binance"
	"log"
)

type (
	AveragePrice map[string]float64
)

func (ap AveragePrice) String() string {
	var buf bytes.Buffer

	for symbol, price := range ap {
		_, err := fmt.Fprintf(&buf, "Average purchase price %v: %.8f\n", symbol, price)
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
		var moneySum, quantitySum, average float64
		for _, point := range tradeList {
			if point.IsBuyer {
				moneySum += point.Quantity * point.Price
				quantitySum += point.Quantity
				average = moneySum / quantitySum
			} else {
				quantitySum -= point.Quantity
				moneySum = average * quantitySum
			}
		}

		prices[symbol] = average
	}

	return prices
}
