package exchangeReq

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeRes"
)

type BinanceBookParams struct {
	Ccy   string `url:"symbol"`
	Limit int    `url:"limit"`
}

func (BinanceBookParams) GetParams(ccy models.Ccy) *models.Request {
	return &models.Request{
		Url:      "https://api.binance.com/api/v3/depth",
		Params:   BinanceBookParams{Ccy: ccy.Currency + ccy.Currency2, Limit: 5},
		Response: &exchangeRes.BinanceBook{},
	}
}
