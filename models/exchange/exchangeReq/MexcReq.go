package exchangeReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes"
)

type MexcBookParams struct {
	Ccy   string `url:"symbol"`
	Limit int    `url:"limit"`
}

func (MexcBookParams) GetParams(ccy models.Ccy) *models.Request {
	return &models.Request{
		Url:      "https://api.mexc.com/api/v3/depth",
		Params:   MexcBookParams{Ccy: ccy.Currency + ccy.Currency2, Limit: 5},
		Response: &exchangeRes.MexcBook{},
	}
}
