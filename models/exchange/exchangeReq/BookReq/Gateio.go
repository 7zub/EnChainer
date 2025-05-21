package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type GateioBookParams struct {
	Ccy string `url:"currency_pair"`
}

func (GateioBookParams) GetParams(ccy any) *models.Request {
	c := ccy.(models.Ccy)

	return &models.Request{
		Url:      "https://api.gateio.ws/api/v4/spot/order_book",
		ReqType:  models.ReqType.Book,
		Params:   GateioBookParams{Ccy: c.Currency + "_" + c.Currency2},
		Response: &BookRes.GateioBook{},
	}
}
