package exchangeReq

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeRes"
)

type BinanceBookParams struct {
	Symbol string
	Limit  int
}

func (BinanceBookParams) GetParams(ccy string) *models.Request {
	return &models.Request{
		Url:      "https://api.binance.com/api/v3/depth?",
		Params:   BinanceBookParams{Symbol: ccy, Limit: 5},
		Response: exchangeRes.BinanceBook{},
	}
}
