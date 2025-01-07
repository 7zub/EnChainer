package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
	"strings"
)

type HuobiBookParams struct {
	Ccy   string `url:"symbol"`
	Limit int    `url:"depth"`
	Type  string `url:"type"`
}

func (HuobiBookParams) GetParams(ccy models.Ccy) *models.Request {
	return &models.Request{
		Url:      "https://api.huobi.pro/market/depth",
		Params:   HuobiBookParams{Ccy: strings.ToLower(ccy.Currency + ccy.Currency2), Limit: 5, Type: "step0"},
		Response: &BookRes.HuobiBook{},
	}
}
