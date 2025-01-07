package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type BinanceBookParams struct {
	Ccy   string `url:"symbol"`
	Limit int    `url:"limit"`
}

func (BinanceBookParams) GetParams(ccy models.Ccy) *models.Request {
	return &models.Request{
		Url:      "https://api.binance.com/api/v3/depth",
		Params:   BinanceBookParams{Ccy: ccy.Currency + ccy.Currency2, Limit: 5},
		Response: &BookRes.BinanceBook{},
	}
}
