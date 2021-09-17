package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/egoreli/binance/pkg/binance"
	"github.com/joho/godotenv"
)

func main() {
	start := time.Now()

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	client, err := binance.NewClient(os.Getenv("API_KEY"), os.Getenv("SECRET_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	balance, err := client.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	printBalance(balance)

	tokens := []string{"KSM", "DOT", "BTC", "ADA", "XRP"}
	pairs := GetTokenPairs(tokens)

	trades, err := client.GetTradeHistory(pairs)
	prices := GetAveragePrices(trades)
	printAveragePrices(prices)

	if err != nil {
		log.Fatal(err)
	}

	rubCourse, err := client.GetRubCourse()
	fmt.Printf("\nCourse: USDTRUB: %v\n", rubCourse)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\nTotal time:", time.Since(start))
}

func printAveragePrices(prices map[string]float64) {
	for symbol, price := range prices {
		fmt.Printf("Average purchase price %v: %.8f\n", symbol, price)
	}
}

func printBalance(balance map[string]float64) {
	fmt.Println("Balance:")
	for token, value := range balance {
		fmt.Printf("%s: %f\n", token, value)
	}
	fmt.Println()
}
