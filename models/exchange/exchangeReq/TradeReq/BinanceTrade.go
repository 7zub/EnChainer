package TradeReq

import (
	"crypto/sha256"
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/BookRes"
	"net/http"
	"strings"
	"time"
)

type BinanceTradeParams struct {
	Ccy    string  `url:"symbol"`
	Side   string  `url:"side"`
	Type   string  `url:"type"`
	Volume float64 `url:"quantity"`
	Price  float64 `url:"price"`
	Live   string  `url:"timeInForce"`
	Time   int64   `url:"timestamp"`
}

func (BinanceTradeParams) GetParams(task any) *models.Request {
	t := task.(models.TradeTask)

	return &models.Request{
		Url:      "https://api.binance.com/api/v3/order",
		ReqType:  "Trade",
		SignType: sha256.New,
		Head: http.Header{
			"X-MBX-APIKEY": []string{models.Conf.ApiKey},
		},
		Params: BinanceTradeParams{
			Ccy:    t.Ccy.Currency + t.Ccy.Currency2,
			Side:   strings.ToUpper(string(t.Stage)),
			Type:   "LIMIT",
			Volume: t.Buy.Volume,
			Price:  t.Buy.Price,
			Live:   "GTC",
			Time:   time.Now().UnixMilli(),
		},
		Response: &BookRes.BinanceBook{},
	}
}
