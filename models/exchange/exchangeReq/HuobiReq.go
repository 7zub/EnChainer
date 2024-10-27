package exchangeReq

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeRes"
	"strings"
)

type HuobiBookParams struct {
	Symbol string
	Depth  int
	Type   string
}

func (HuobiBookParams) GetParams(ccy models.Ccy) *models.Request {
	return &models.Request{
		Url:      "https://api.huobi.pro/market/depth?",
		Params:   HuobiBookParams{Symbol: strings.ToLower(ccy.Currency + ccy.Currency2), Depth: 5, Type: "step0"},
		Response: &exchangeRes.HuobiBook{},
	}
}
