package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/joho/godotenv"
)

type trade struct {
	price      float64
	quantity   float64
	commission float64
	isBuyer    bool
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic("No .env file found")
	}

	start := time.Now()

	client := binance.NewClient(os.Getenv("API_KEY"), os.Getenv("SECRET_KEY"))

	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	balance := make(map[string]float64)
	for _, pos := range res.Balances {
		free, _ := strconv.ParseFloat(pos.Free, 64)
		if free != 0 {
			balance[pos.Asset] = free
		}
	}
	fmt.Printf("%v\n\n", balance)

	var pairs []string
	for coin := range balance {
		if coin != "USDT" {
			pairs = append(pairs, coin+"USDT")
		}
	}
	//pairs = append(pairs, "KSMUSDT")
	fmt.Printf("%v\n\n", pairs)

	trades := make(map[string][]trade)
	for _, pair := range pairs {
		tradesList, err := client.NewListTradesService().Symbol(pair).
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, tradePoint := range tradesList {
			price, _ := strconv.ParseFloat(tradePoint.Price, 64)
			quantity, _ := strconv.ParseFloat(tradePoint.Quantity, 64)
			commission, _ := strconv.ParseFloat(tradePoint.Commission, 64)
			isBuyer := tradePoint.IsBuyer

			trades[tradePoint.Symbol] = append(trades[tradePoint.Symbol],
				trade{
					price:      price,
					quantity:   quantity,
					commission: commission,
					isBuyer:    isBuyer,
				})
		}
	}

	for symbol, tradeList := range trades {
		var moneySum, quantitySum, average float64
		for _, point := range tradeList {
			if point.isBuyer {
				moneySum += point.quantity * point.price
				quantitySum += point.quantity
				average = moneySum / quantitySum
			} else {
				quantitySum -= point.quantity
				moneySum = average * quantitySum
			}
		}

		fmt.Printf("Average purchase price %v: %.8f\n", symbol, average)
	}

	price, err := client.NewListPricesService().Symbol("USDTRUB").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("\n%s: %v\n", price[0].Symbol, price[0].Price)

	fmt.Println("\n\nTotal time:", time.Since(start))
}
