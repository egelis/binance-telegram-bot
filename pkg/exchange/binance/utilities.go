package binance

import (
	binanceAPI "github.com/adshao/go-binance/v2"
	"github.com/egelis/binance/pkg/exchange"
	"github.com/shopspring/decimal"
)

func getTradePoint(tpv3 *binanceAPI.TradeV3) (exchange.TradePoint, error) {
	price, err := decimal.NewFromString(tpv3.Price)
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
		Price:           price,
		Quantity:        quantity,
		Commission:      commission,
		CommissionAsset: commissionAsset,
		IsBuyer:         isBuyer,
	}, nil
}
