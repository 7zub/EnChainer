package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type MexcBookParams struct {
	Ccy   string `url:"symbol"`
	Limit int    `url:"limit"`
}

func (MexcBookParams) GetParams(pair any) *models.Request {
	p := pair.(*models.TradePair)

	var url, ccy string
	switch p.Market {
	case models.Market.Spot:
		url = "https://api.mexc.com/api/v3/depth"
		ccy = p.Ccy.Currency + p.Ccy.Currency2
	case models.Market.Futures:
		url = "https://contract.mexc.com/api/v1/contract/depth/" + p.Ccy.Currency + "_" + p.Ccy.Currency2
	}

	return &models.Request{
		Url:      url,
		ReqType:  models.ReqType.Book,
		Params:   MexcBookParams{Ccy: ccy, Limit: 5},
		Response: &BookRes.MexcBook{},
	}
}
