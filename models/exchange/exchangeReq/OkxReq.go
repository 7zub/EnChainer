package exchangeReq

import (
	"awesomeProject/models"
)

type OkxBookParams struct {
	Ccy   string `url:"instId"`
	Limit int    `url:"sz"`
}

func (OkxBookParams) GetParams(ccy models.Ccy) *models.Request {
	return &models.Request{
		Url:    "https://www.okx.com/api/v5/market/books?",
		Params: OkxBookParams{Ccy: ccy.Currency + "-" + ccy.Currency2, Limit: 5},
		//Response: &exchangeRes.HuobiBook{},
	}
}
