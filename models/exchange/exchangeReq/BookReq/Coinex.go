package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type CoinexBookParams struct {
	Ccy     string `url:"market"`
	Limit   int    `url:"limit"`
	Decimal string `url:"interval"`
}

func (CoinexBookParams) GetParams(pair any) *models.Request {
	p := pair.(*models.TradePair)

	var mark string
	switch p.Market {
	case models.Market.Spot:
		mark = "spot"
	case models.Market.Feature:
		mark = "futures"
	}

	return &models.Request{
		Url:      "https://api.coinex.com/v2/" + mark + "/depth",
		ReqType:  models.ReqType.Book,
		Params:   CoinexBookParams{Ccy: p.Ccy.Currency + p.Ccy.Currency2, Limit: 5, Decimal: "0"},
		Response: &BookRes.CoinexBook{},
	}
}
