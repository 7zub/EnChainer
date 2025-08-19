package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type OkxBookParams struct {
	Ccy   string `url:"instId"`
	Limit int    `url:"sz"`
}

func (OkxBookParams) GetParams(pair any) *models.Request {
	p := pair.(*models.TradePair)

	var ccy string
	switch p.Market {
	case models.Market.Spot:
		ccy = p.Ccy.Currency + "-" + p.Ccy.Currency2
	case models.Market.Features:
		ccy = p.Ccy.Currency + "-" + p.Ccy.Currency2 + "-SWAP"
	}

	return &models.Request{
		Url:      "https://okx.com/api/v5/market/books",
		ReqType:  models.ReqType.Book,
		Params:   OkxBookParams{Ccy: ccy, Limit: 5},
		Response: &BookRes.OkxBook{},
	}
}
