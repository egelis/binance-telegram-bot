package main

import (
	"fmt"
	"log"
	"os"

	"github.com/egelis/binance/pkg/exchange/binance"
	"github.com/egelis/binance/pkg/portfolio"
	"github.com/joho/godotenv"
)

var TOKENS = []string{"KSM", "DOT", "BTC", "ADA", "XRP", "ATOM"}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	binanceExchange, err := binance.NewClient(os.Getenv("API_KEY"), os.Getenv("SECRET_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	binancePortfolio := portfolio.NewPortfolio(binanceExchange)

	fmt.Println(binancePortfolio.GetTradeHistoryForPair("BTCUSDT"))

	fmt.Println()

	for _, token := range TOKENS {
		fmt.Println(binancePortfolio.GetTokenStatistic(token))
	}
}
