package binance

import (
	binanceAPI "github.com/adshao/go-binance/v2"
	"github.com/egelis/binance/pkg/exchange"
	"github.com/shopspring/decimal"
	"time"
)

func getTradePoint(tpv3 *binanceAPI.TradeV3) (exchange.TradePoint, error) {
	t := time.Unix(tpv3.Time/1000, 0)

	price, err := decimal.NewFromString(tpv3.Price)
	if err != nil {
		return exchange.TradePoint{}, err
	}

	purchasePrice, err := decimal.NewFromString(tpv3.QuoteQuantity)
	if err != nil {
		return exchange.TradePoint{}, err
	}

	quantity, err := decimal.NewFromString(tpv3.Quantity)
	if err != nil {
		return exchange.TradePoint{}, err
	}

	commission, err := decimal.NewFromString(tpv3.Commission)
	if err != nil {
		return exchange.TradePoint{}, err
	}

	commissionAsset := tpv3.CommissionAsset

	isBuyer := tpv3.IsBuyer

	return exchange.TradePoint{
		Time:            t,
		Price:           price,
		PurchasePrice:   purchasePrice,
		Quantity:        quantity,
		Commission:      commission,
		CommissionAsset: commissionAsset,
		IsBuyer:         isBuyer,
	}, nil
}
