package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type CoinexBookParams struct {
	Ccy     string  `url:"market"`
	Limit   int     `url:"limit"`
	Decimal float64 `url:"merge"`
}

func (CoinexBookParams) GetParams(ccy any) *models.Request {
	c := ccy.(models.Ccy)

	return &models.Request{
		Url:      "https://api.coinex.com/v1/market/depth",
		ReqType:  "Book",
		Params:   CoinexBookParams{Ccy: c.Currency + c.Currency2, Limit: 5, Decimal: 0.000000000001},
		Response: &BookRes.CoinexBook{},
	}
}
