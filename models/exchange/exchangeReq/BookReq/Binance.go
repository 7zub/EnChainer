package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type BinanceBookParams struct {
	Ccy   string `url:"symbol"`
	Limit int    `url:"limit"`
}

func (BinanceBookParams) GetParams(ccy any) *models.Request {
	c := ccy.(models.Ccy)

	return &models.Request{
		Url:      "https://api.binance.com/api/v3/depth",
		ReqType:  "Book",
		Params:   BinanceBookParams{Ccy: c.Currency + c.Currency2, Limit: 5},
		Response: &BookRes.BinanceBook{},
	}
}
