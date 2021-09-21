package main

import (
	"fmt"
	"github.com/egelis/binance/pkg/exchange/binance"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	TOKENS = []string{"KSM", "DOT", "BTC", "ADA", "XRP", "ATOM"}
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

	GetTokensStatistic(client)

	fmt.Println("\n\nTotal time:", time.Since(start))
}
