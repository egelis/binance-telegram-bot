package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/egelis/binance/pkg/binance"
	"github.com/joho/godotenv"
)

var (
	TOKENS = []string{"KSM", "DOT", "BTC", "ADA", "XRP", "ATOM"}
)

func getStakingAmount(tradeHistory []binance.TradePoint, divSum float64, balance float64) float64 {
	var amount float64
	for _, tradePoint := range tradeHistory {
		if tradePoint.IsBuyer {
			amount += tradePoint.Quantity
			amount -= tradePoint.Commission
		} else {
			amount -= tradePoint.Quantity
		}
	}

	amount += divSum
	amount -= balance

	return amount
}

func printTokensStatistic(c *binance.Client) {
	balance, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}

	pairs := GetTokenPairs(TOKENS)
	tradeHistory, err := c.GetTradeHistory(pairs)
	if err != nil {
		log.Fatal(err)
	}

	averagePrice := GetAveragePrices(tradeHistory)

	for _, token := range TOKENS {
		dividends, err := c.GetTokenDividends(token)
		if err != nil {
			log.Fatal(err)
		}

		var divSum float64
		for _, div := range dividends {
			divSum += div
		}

		staking := getStakingAmount(tradeHistory[token+"USDT"], divSum, balance[token])

		fmt.Printf("Token: %s\n"+
			"Balance: %f\n"+
			"In staking: %.8f\n"+
			"Average: %.2f\n"+
			"Dividends: %f\n\n",
			token, balance[token], staking, averagePrice[token+"USDT"], divSum)
	}
}

func main() {
	start := time.Now()

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	client, err := binance.NewClient(os.Getenv("API_KEY"), os.Getenv("SECRET_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	printTokensStatistic(client)

	//client.GetTokenDividends()
	//
	//balance, err := client.GetBalance()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(balance)
	//
	//pairs := GetTokenPairs(TOKENS)
	//
	//trades, err := client.GetTradeHistory(pairs)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//averagePrices := GetAveragePrices(trades)
	//fmt.Println(averagePrices)
	//
	//rubCourse, err := client.GetRubCourse()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("USDTRUB: %v\n", rubCourse)

	fmt.Println("\n\nTotal time:", time.Since(start))
}
