package BookReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
)

type MexcBookParams struct {
	Ccy   string `url:"symbol"`
	Limit int    `url:"limit"`
}

func (MexcBookParams) GetParams(ccy any) *models.Request {
	c := ccy.(models.Ccy)

	return &models.Request{
		Url:      "https://api.mexc.com/api/v3/depth",
		ReqType:  models.ReqType.Book,
		Params:   MexcBookParams{Ccy: c.Currency + c.Currency2, Limit: 5},
		Response: &BookRes.MexcBook{},
	}
}
