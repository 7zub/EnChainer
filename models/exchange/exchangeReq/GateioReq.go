package exchangeReq

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeRes"
)

type GateioBookParams struct {
	Currency_pair string
}

func (GateioBookParams) GetParams(ccy string) *models.Request {
	return &models.Request{
		Url:      "https://api.binance.com/api/v3/depth?",
		Params:   GateioBookParams{Currency_pair: ccy},
		Response: exchangeRes.GateioBook{},
	}
}
