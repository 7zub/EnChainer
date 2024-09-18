package exchangeReq

import "awesomeProject/models"

type GateioBookParams struct {
	Currency_pair string
}

func (GateioBookParams) GetReqBook() models.Request {
	return models.Request{
		Url:      "https://api.gateio.ws/api/v4/spot/order_book?",
		Currency: "NEO",
		Params:   GateioBookParams{Currency_pair: "NEO"},
	}
}
