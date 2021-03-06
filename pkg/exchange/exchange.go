package exchange

import (
	"github.com/shopspring/decimal"
	"time"
)

type Balance map[string]decimal.Decimal

type Dividends []decimal.Decimal

type TradePoint struct {
	Time            time.Time
	Price           decimal.Decimal
	PurchasePrice   decimal.Decimal
	Quantity        decimal.Decimal
	Commission      decimal.Decimal
	CommissionAsset string
	IsBuyer         bool
}

type Exchange interface {
	GetBalance() (Balance, error)
	GetTokenPrice(token string) (decimal.Decimal, error)
	GetPairTradeHistory(pair string) ([]TradePoint, error)
	GetTokenDividends(token string) (Dividends, error)
}
