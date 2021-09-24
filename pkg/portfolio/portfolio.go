package portfolio

import (
	"bytes"
	"fmt"
	"log"

	"github.com/egelis/binance/pkg/exchange"
	"github.com/shopspring/decimal"
)

type TokenStatistic struct {
	Balance      decimal.Decimal
	InStaking    decimal.Decimal
	AveragePrice decimal.Decimal
	Dividends    decimal.Decimal
	CurrentPrice decimal.Decimal
	Profit       decimal.Decimal
}

type Portfolio struct {
	exchange exchange.Exchange
}

func NewPortfolio(exchange exchange.Exchange) *Portfolio {
	return &Portfolio{
		exchange: exchange,
	}
}

func (p *Portfolio) GetTokenStatistic(token string) string {
	stat, err := p.getTokenStatistic(token)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	_, err = fmt.Fprintf(&buf,
		"Токен:            %s\n"+
			"Баланс:           %v\n"+
			"В стейкинге:      %v\n"+
			"Ср. цена покупки: %v USDT\n"+
			"Дивиденды:        %v\n"+
			"Текущая цена:     %v\n"+
			"Профит:           %s (%s)\n",
		token,
		stat.Balance.StringFixed(8),
		stat.InStaking.StringFixed(8),
		stat.AveragePrice.StringFixed(2),
		stat.Dividends.StringFixed(8),
		stat.CurrentPrice.StringFixed(2),
		getProfitPercents(stat.Profit),
		getProfitUSD(stat.Balance, stat.InStaking, stat.AveragePrice, stat.CurrentPrice))

	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
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
