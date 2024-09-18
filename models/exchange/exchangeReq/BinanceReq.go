package exchangeReq

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeRes"
)

type BinanceBookParams struct {
	Symbol string
	Limit  int
}

func (BinanceBookParams) GetParams() models.Request {
	return models.Request{
		Url:      "https://api.binance.com/api/v3/depth?",
		Params:   BinanceBookParams{Symbol: "NEO", Limit: 5},
		Response: exchangeRes.BinanceBook{},
	}
}
