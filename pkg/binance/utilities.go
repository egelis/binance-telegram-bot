package binance

import (
	binanceAPI "github.com/adshao/go-binance/v2"
	"strconv"
)

func getTradePoint(tpv3 *binanceAPI.TradeV3) (TradePoint, error) {
	price, err := strconv.ParseFloat(tpv3.Price, 64)
	if err != nil {
		return TradePoint{}, err
	}

	quantity, err := strconv.ParseFloat(tpv3.Quantity, 64)
	if err != nil {
		return TradePoint{}, err
	}

	commission, err := strconv.ParseFloat(tpv3.Commission, 64)
	if err != nil {
		return TradePoint{}, err
	}

	commissionAsset := tpv3.CommissionAsset

	isBuyer := tpv3.IsBuyer

	return TradePoint{
		Price:           price,
		Quantity:        quantity,
		Commission:      commission,
		CommissionAsset: commissionAsset,
		IsBuyer:         isBuyer,
	}, nil
}
