package TradeReq

import (
	"bytes"
	"crypto/sha256"
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/TradeRes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OkxTradeParams struct {
	Ccy       string  `url:"-" json:"instId"`
	Side      string  `url:"-" json:"side"`
	Type      string  `url:"-" json:"ordType"`
	Volume    float64 `url:"-" json:"sz"`
	Price     float64 `url:"-" json:"px"`
	Margin    string  `url:"-" json:"tdMode"`
	MarginCcy string  `url:"-" json:"ccy"`
}

func (OkxTradeParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)
	endpoint := "/api/v5/trade/order"

	return &models.Request{
		Url:     "https://www.okx.com" + endpoint,
		ReqType: models.ReqType.Trade,
		SignWay: func(rq *http.Request) {

			jsonBody, _ := json.Marshal(OkxTradeParams{
				Ccy:       t.Currency + "-" + t.Currency2,
				Side:      strings.ToLower(string(t.Side)),
				Type:      "limit",
				Volume:    t.Volume,
				Price:     t.Price,
				Margin:    "cross",
				MarginCcy: t.Currency2,
			})

			timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
			payload := fmt.Sprintf("%sPOST%s%s", timestamp, endpoint, string(jsonBody[:]))
			sign := models.Sign(payload, models.Conf.Exchanges[string(t.Ex)].SecretKey, sha256.New, "base64")

			rq.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
			rq.Header.Add("Content-Type", "application/json")
			rq.Header.Add("OK-ACCESS-KEY", models.Conf.Exchanges[string(t.Ex)].ApiKey)
			rq.Header.Add("OK-ACCESS-SIGN", sign)
			rq.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
			rq.Header.Add("OK-ACCESS-PASSPHRASE", models.Conf.Exchanges[string(t.Ex)].PassPhrase)

		},
		Params:   OkxTradeParams{},
		Response: &TradeRes.OkxTrade{},
	}
}
