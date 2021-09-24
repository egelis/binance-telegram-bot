package main

import (
	"github.com/egelis/binance/pkg/portfolio"
	"log"
	"os"

	"github.com/egelis/binance/pkg/exchange/binance"
	"github.com/joho/godotenv"
)

var (
	TOKENS = []string{"KSM", "DOT", "BTC", "ADA", "XRP", "ATOM"}
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	binanceExchange, err := binance.NewClient(os.Getenv("API_KEY"), os.Getenv("SECRET_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	binancePortfolio := portfolio.NewPortfolio(binanceExchange, TOKENS)

	//fmt.Println(binancePortfolio.GetTradeHistoryForPair("DOTUSDT"))
	binancePortfolio.GetTokenListStatistic()
}
