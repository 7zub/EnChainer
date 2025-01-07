package BookReq

import (
	"time"
)

type BinanceTradeParams struct {
	Ccy    string    `url:"symbol"`
	Side   string    `url:"side"`
	Type   string    `url:"type"`
	Volume float64   `url:"quantity"`
	Price  float64   `url:"price"`
	Live   string    `url:"timeInForce"`
	Time   time.Time `url:"timestamp"`
}

/*func (BinanceTradeParams) GetParams(ccy models.Ccy, side string, volume float64, price float64) *models.Request {
	return &models.Request{
		Url:      "https://api.binance.com/api/v3/order",
		Params:   BinanceTradeParams{Ccy: ccy.Currency + ccy.Currency2, Type: "LIMIT"},
		Response: &BookRes.BinanceBook{},
	}
}*/
