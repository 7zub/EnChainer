package exchangeReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes"
)

type KucoinBookParams struct {
	Ccy string `url:"symbol"`
}

func (KucoinBookParams) GetParams(ccy models.Ccy) *models.Request {
	return &models.Request{
		Url:      "https://api.kucoin.com/api/v1/market/orderbook/level2_20",
		Params:   KucoinBookParams{Ccy: ccy.Currency + "-" + ccy.Currency2},
		Response: &exchangeRes.KucoinBook{},
	}
}
