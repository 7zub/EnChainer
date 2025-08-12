package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type KucoinBookParams struct {
	Ccy string `url:"symbol"`
}

func (KucoinBookParams) GetParams(pair any) *models.Request {
	p := pair.(*models.TradePair)

	var url, ccy string
	switch p.Market {
	case models.Market.Spot:
		url = "https://api.kucoin.com/api/v1/market/orderbook/level2_20"
		ccy = p.Ccy.Currency + "-" + p.Ccy.Currency2

	case models.Market.Feature:
		url = "https://api-futures.kucoin.com/api/v1/level2/depth20"
		ccy = p.Ccy.Currency + p.Ccy.Currency2 + "M"
	}

	return &models.Request{
		Url:      url,
		ReqType:  models.ReqType.Book,
		Params:   KucoinBookParams{Ccy: ccy},
		Response: &BookRes.KucoinBook{},
	}
}
