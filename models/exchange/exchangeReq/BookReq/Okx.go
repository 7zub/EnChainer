package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type OkxBookParams struct {
	Ccy   string `url:"instId"`
	Limit int    `url:"sz"`
}

func (OkxBookParams) GetParams(ccy any) *models.Request {
	c := ccy.(models.Ccy)

	return &models.Request{
		Url:      "https://www.okx.com/api/v5/market/books",
		Params:   OkxBookParams{Ccy: c.Currency + "-" + c.Currency2, Limit: 5},
		Response: &BookRes.OkxBook{},
	}
}
