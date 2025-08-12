package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
	"strings"
)

type GateioBookParams struct {
	Ccy      string `url:"currency_pair"`
	Contract string `url:"contract"`
}

func (GateioBookParams) GetParams(pair any) *models.Request {
	p := pair.(*models.TradePair)

	var mark string
	var params GateioBookParams
	switch p.Market {
	case models.Market.Spot:
		mark = "spot"
		params = GateioBookParams{Ccy: p.Ccy.Currency + "_" + p.Ccy.Currency2}
	case models.Market.Feature:
		mark = "futures/" + strings.ToLower(p.Ccy.Currency2)
		params = GateioBookParams{Contract: p.Ccy.Currency + "_" + p.Ccy.Currency2}
	}

	return &models.Request{
		Url:      "https://api.gateio.ws/api/v4/" + mark + "/order_book",
		ReqType:  models.ReqType.Book,
		Params:   params,
		Response: &BookRes.GateioBook{},
	}
}
