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
	fmt.Println(balance)

	tokens := []string{"KSM", "DOT", "BTC", "ADA", "XRP"}
	pairs := GetTokenPairs(tokens)

	trades, err := client.GetTradeHistory(pairs)
	if err != nil {
		log.Fatal(err)
	}

	averagePrices := GetAveragePrices(trades)
	fmt.Println(averagePrices)

	rubCourse, err := client.GetRubCourse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("USDTRUB: %v\n", rubCourse)

	fmt.Println("\n\nTotal time:", time.Since(start))
}
