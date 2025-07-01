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

func (BybitBookParams) GetParams(ccy any) *models.Request {
	c := ccy.(models.Ccy)

	return &models.Request{
		Url:      "https://api.bybit.com/v5/market/orderbook",
		ReqType:  models.ReqType.Book,
		Params:   BybitBookParams{Ccy: c.Currency + c.Currency2, Category: "spot", Limit: 5},
		Response: &BookRes.BybitBook{},
	}
}
