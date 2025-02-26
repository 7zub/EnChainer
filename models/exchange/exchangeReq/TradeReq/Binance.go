package TradeReq

import (
	"crypto/sha256"
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/TradeRes"
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
		Url:     "https://api.binance.com/api/v3/order",
		ReqType: "Trade",
		SignWay: func(rq *http.Request) {
			rq.Header.Add("X-MBX-APIKEY", models.Conf.Exchanges[t.Buy.Ex].ApiKey)
			q := rq.URL.Query()
			sign := models.Sign(rq.URL.Query().Encode(), models.Conf.Exchanges[t.Buy.Ex].SecretKey, sha256.New)
			q.Add("signature", sign)
			rq.URL.RawQuery = q.Encode()
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
		Response: &TradeRes.BinanceTrade{},
	}
}
