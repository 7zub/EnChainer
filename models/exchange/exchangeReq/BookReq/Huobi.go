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

func (HuobiBookParams) GetParams(ccy any) *models.Request {
	c := ccy.(models.Ccy)

	return &models.Request{
		Url:      "https://api.huobi.pro/market/depth",
		ReqType:  models.ReqType.Book,
		Params:   HuobiBookParams{Ccy: strings.ToLower(c.Currency + c.Currency2), Limit: 5, Type: "step0"},
		Response: &BookRes.HuobiBook{},
	}
}
