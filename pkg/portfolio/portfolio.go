package portfolio

import "github.com/egelis/binance/pkg/exchange"

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
