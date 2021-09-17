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
		log.Panic("No .env file found")
	}

	client := binance.NewClient(os.Getenv("API_KEY"), os.Getenv("SECRET_KEY"))

	balance, _ := client.GetBalance()
	fmt.Printf("%v\n\n", balance)

	pairs := client.MakeTokenPairsFromBalance(balance)
	fmt.Printf("%v\n\n", pairs)

	trades, _ := client.GetTradeHistory(pairs)
	client.PrintAveragePrices(trades)

	rubCourse, _ := client.GetRubCourse()
	fmt.Printf("\nUSDTRUB: %v\n", rubCourse)

	fmt.Println("\n\nTotal time:", time.Since(start))
}
