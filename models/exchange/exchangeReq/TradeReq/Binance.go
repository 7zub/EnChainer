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
	Ccy        string  `url:"symbol"`
	Side       string  `url:"side"`
	Type       string  `url:"type"`
	Volume     float64 `url:"quantity"`
	Price      float64 `url:"price"`
	Live       string  `url:"timeInForce"`
	Time       int64   `url:"timestamp"`
	AutoBorrow string  `url:"sideEffectType"`
}

func (BinanceTradeParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)

	return &models.Request{
		Url:     "https://api.binance.com/sapi/v1/margin/order",
		ReqType: models.ReqType.Trade,
		SignWay: func(rq *http.Request) {
			rq.Header.Add("X-MBX-APIKEY", models.Conf.Exchanges[string(t.Ex)].ApiKey)
			q := rq.URL.Query()
			sign := models.Sign(rq.URL.Query().Encode(), models.Conf.Exchanges[string(t.Ex)].SecretKey, sha256.New, "hex")
			q.Add("signature", sign)
			rq.URL.RawQuery = q.Encode()
		},
		Params: BinanceTradeParams{
			Ccy:        t.Currency + t.Currency2,
			Side:       strings.ToUpper(string(t.Side)),
			Type:       "LIMIT",
			Volume:     t.Volume,
			Price:      t.Price,
			Live:       "GTC",
			Time:       time.Now().UnixMilli(),
			AutoBorrow: "AUTO_BORROW_REPAY",
		},
		Response: &TradeRes.BinanceTrade{},
	}
}
