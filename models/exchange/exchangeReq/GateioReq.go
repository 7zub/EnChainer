package exchangeReq

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeRes"
)

type GateioBookParams struct {
	Ccy string `url:"currency_pair"`
}

func (GateioBookParams) GetParams(ccy models.Ccy) *models.Request {
	return &models.Request{
		Url:      "https://api.gateio.ws/api/v4/spot/order_book",
		Params:   GateioBookParams{Ccy: ccy.Currency + "_" + ccy.Currency2},
		Response: &exchangeRes.GateioBook{},
	}
}
