package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
	"strings"
)

type HuobiBookParams struct {
	Ccy      string `url:"symbol"`
	Contract string `url:"contract_code"`
	Limit    string `url:"depth"`
	Type     string `url:"type"`
}

func (HuobiBookParams) GetParams(pair any) *models.Request {
	p := pair.(*models.TradePair)

	var url string
	var params HuobiBookParams
	switch p.Market {
	case models.Market.Spot:
		url = "https://api.huobi.pro/market/depth"
		params = HuobiBookParams{Ccy: strings.ToLower(p.Ccy.Currency + p.Ccy.Currency2), Limit: "5", Type: "step0"}
	case models.Market.Features:
		url = "https://api.hbdm.com/linear-swap-ex/market/depth"
		params = HuobiBookParams{Contract: strings.ToUpper(p.Ccy.Currency + "-" + p.Ccy.Currency2), Type: "step0"}
	}

	return &models.Request{
		Url:      url,
		ReqType:  models.ReqType.Book,
		Params:   params,
		Response: &BookRes.HuobiBook{},
	}
}
