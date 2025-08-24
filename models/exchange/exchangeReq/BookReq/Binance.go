package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type BinanceBookParams struct {
	Ccy   string `url:"symbol"`
	Limit int    `url:"limit"`
}

func (BinanceBookParams) GetParams(pair any) *models.Request {
	p := pair.(*models.TradePair)

	var url string
	switch p.Market {
	case models.Market.Spot:
		url = "https://api.binance.com/api/v3/depth"
	case models.Market.Futures:
		url = "https://fapi.binance.com/fapi/v1/depth"
	}

	return &models.Request{
		Url:      url,
		ReqType:  models.ReqType.Book,
		Params:   BinanceBookParams{Ccy: p.Ccy.Currency + p.Ccy.Currency2, Limit: 5},
		Response: &BookRes.BinanceBook{},
	}
}
