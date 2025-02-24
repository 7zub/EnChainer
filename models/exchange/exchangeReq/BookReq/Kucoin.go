package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type KucoinBookParams struct {
	Ccy string `url:"symbol"`
}

func (KucoinBookParams) GetParams(ccy any) *models.Request {
	c := ccy.(models.Ccy)

	return &models.Request{
		Url:      "https://api.kucoin.com/api/v1/market/orderbook/level2_20",
		ReqType:  "Book",
		Params:   KucoinBookParams{Ccy: c.Currency + "-" + c.Currency2},
		Response: &BookRes.KucoinBook{},
	}
}
