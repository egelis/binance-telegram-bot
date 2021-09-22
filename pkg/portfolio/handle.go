package portfolio

import (
	"bytes"
	"fmt"
	"log"
)

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
