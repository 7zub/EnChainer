package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type BybitBookParams struct {
	Ccy      string `url:"symbol"`
	Limit    int    `url:"limit"`
	Category string `url:"category"`
}

func (BybitBookParams) GetParams(pair any) *models.Request {
	p := pair.(*models.TradePair)

	var mark string
	switch p.Market {
	case models.Market.Spot:
		mark = "spot"
	case models.Market.Futures:
		mark = "linear"
	}

	return &models.Request{
		Url:      "https://api.bybit.com/v5/market/orderbook",
		ReqType:  models.ReqType.Book,
		Params:   BybitBookParams{Ccy: p.Ccy.Currency + p.Ccy.Currency2, Category: mark, Limit: 5},
		Response: &BookRes.BybitBook{},
	}
}
