package TradeReq

import (
	"bytes"
	"crypto/sha256"
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/TradeRes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type KucoinTradeParams struct {
	Id     string  `url:"-" json:"clientOid"`
	Ccy    string  `url:"-" json:"symbol"`
	Side   string  `url:"-" json:"side"`
	Type   string  `url:"-" json:"type"`
	Volume float64 `url:"-" json:"size"`
	Price  float64 `url:"-" json:"price"`
	Margin string  `url:"-" json:"marginMode"`
}

func (KucoinTradeParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)
	endpoint := "/api/v1/orders"

	var url, ccy string
	switch t.Market {
	case models.Market.Spot:
		url = "https://api.kucoin.com"
		ccy = t.Ccy.Currency + "-" + t.Ccy.Currency2
	case models.Market.Futures:
		url = "https://api-futures.kucoin.com"
		ccy = t.Ccy.Currency + t.Ccy.Currency2 + "M"
	}

	return &models.Request{
		Url:     url + endpoint,
		ReqType: models.ReqType.Trade,
		SignWay: func(rq *http.Request) {
			jsonBody, _ := json.Marshal(KucoinTradeParams{
				Id:     fmt.Sprintf("kc%06d", rand.Intn(1000000)),
				Ccy:    ccy,
				Side:   string(t.Side),
				Type:   "limit",
				Volume: t.Volume,
				Price:  t.Price,
				Margin: "CROSS",
			})

			timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
			payload := fmt.Sprintf("%sPOST%s%s", timestamp, endpoint, string(jsonBody[:]))
			sign := models.Sign(payload, models.Conf.Exchanges[string(t.Ex)].SecretKey, sha256.New, "base64")
			passphrase := models.Sign(models.Conf.Exchanges[string(t.Ex)].PassPhrase, models.Conf.Exchanges[string(t.Ex)].SecretKey, sha256.New, "base64")

			rq.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
			rq.Header.Add("Content-Type", "application/json")
			rq.Header.Add("KC-API-KEY", models.Conf.Exchanges[string(t.Ex)].ApiKey)
			rq.Header.Add("KC-API-SIGN", sign)
			rq.Header.Add("KC-API-TIMESTAMP", timestamp)
			rq.Header.Add("KC-API-PASSPHRASE", passphrase)
			rq.Header.Set("KC-API-KEY-VERSION", "3")
		},
		Params:   KucoinTradeParams{},
		Response: &TradeRes.KucoinTrade{},
	}
}
