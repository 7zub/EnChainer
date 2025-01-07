package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type BybitBookParams struct {
	Ccy      string `url:"symbol"`
	Limit    int    `url:"limit"`
	Category string `url:"category"`
}

func (BybitBookParams) GetParams(ccy models.Ccy) *models.Request {
	return &models.Request{
		Url:      "https://api.bybit.com/v5/market/orderbook",
		Params:   BybitBookParams{Ccy: ccy.Currency + ccy.Currency2, Category: "linear", Limit: 5},
		Response: &BookRes.BybitBook{},
	}
}
