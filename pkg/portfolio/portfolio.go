package portfolio

import (
	"bytes"
	"fmt"
	"github.com/egelis/binance/pkg/exchange"
	"github.com/shopspring/decimal"
	"log"
)

type Portfolio struct {
	exchange exchange.Exchange

	tokens []string
}

func NewPortfolio(exchange exchange.Exchange, tokens []string) *Portfolio {
	return &Portfolio{
		exchange: exchange,
		tokens:   tokens,
	}
}

func (p *Portfolio) GetTradeHistoryForPair(pair string) string {
	tradePoints, err := p.exchange.GetPairTradeHistory(pair)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	_, err = fmt.Fprintf(&buf, "%s\n"+
		"-------------------------------------------",
		pair)

	if err != nil {
		log.Fatal(err)
	}

	for _, point := range tradePoints {

		_, err = fmt.Fprintf(&buf,
			"\nДата покупки: %v\n"+
				"Цена актива:  %v\n"+
				"Цена покупки: %v\n"+
				"Количество:   %v\n"+
				"Комиссия:     %v %v\n"+
				"-------------------------------------------",
			point.Time,
			point.Price,
			point.PurchasePrice,
			point.Quantity,
			point.Commission,
			point.CommissionAsset)

		if err != nil {
			log.Fatal(err)
		}
	}

	return buf.String()
}

func (p *Portfolio) GetTokenListStatistic() {
	balance, err := p.exchange.GetBalance()
	if err != nil {
		log.Fatal(err)
	}

	for _, token := range p.tokens {
		dividends, err := p.exchange.GetTokenDividends(token)
		if err != nil {
			log.Fatal(err)
		}

		divSum := decimal.NewFromFloat(0)
		for _, div := range dividends {
			divSum = divSum.Add(div)
		}

		tokenUSDT := token + "USDT"

		tradeHistory, err := p.exchange.GetPairTradeHistory(tokenUSDT)
		if err != nil {
			log.Fatal(err)
		}

		staking := getStakingAmount(tradeHistory, divSum, balance[token])

		currentPrice, err := p.exchange.GetTokenPrice(tokenUSDT)
		if err != nil {
			log.Fatal(err)
		}

		averagePrice := getAveragePrices(tradeHistory)
		var profit decimal.Decimal
		if !averagePrice.Equal(decimal.NewFromFloat(0)) {
			profit = currentPrice.Div(averagePrice)
		} else {
			profit = decimal.NewFromFloat(0)
		}

		fmt.Printf("Token:         %s\n"+
			"Balance:       %v\n"+
			"In staking:    %v\n"+
			"Average:       %v\n"+
			"Dividends:     %v\n"+
			"Current price: %v\n"+
			"Profit:        %v\n\n",
			token,
			balance[token].StringFixed(8),
			staking.StringFixed(8),
			averagePrice.StringFixed(2),
			divSum.StringFixed(8),
			currentPrice.StringFixed(2),
			profit.StringFixed(4))
	}
}
