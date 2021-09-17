package main

import (
	"github.com/egoreli/binance/pkg/binance"
)

func GetTokenPairs(tokens []string) []string {
	var pairs []string

	for _, coin := range tokens {
		if coin != "USDT" {
			pairs = append(pairs, coin+"USDT")
		}
	}

	return pairs
}

func GetAveragePrices(trades map[string][]binance.TradePoint) map[string]float64 {
	prices := make(map[string]float64)

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
