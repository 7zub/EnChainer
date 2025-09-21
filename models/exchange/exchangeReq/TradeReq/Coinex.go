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
	"strconv"
	"strings"
	"time"
)

type CoinexTradeParams struct {
	Ccy    string  `url:"-" json:"market"`
	Side   string  `url:"-" json:"side"`
	Type   string  `url:"-" json:"type"`
	Volume float64 `url:"-" json:"amount"`
	Price  float64 `url:"-" json:"price"`
	Margin string  `url:"-" json:"market_type"`
}

func (CoinexTradeParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)

	var endpoint string
	switch t.Market {
	case models.Market.Spot:
		endpoint = "/v2/spot/order"
	case models.Market.Futures:
		endpoint = "/v2/futures/order"
	}

	return &models.Request{
		Url:     "https://api.coinex.com" + endpoint,
		ReqType: models.ReqType.Trade,
		SignWay: func(rq *http.Request) {
			jsonBody, _ := json.Marshal(CoinexTradeParams{
				Ccy:    t.Currency + t.Currency2,
				Side:   strings.ToLower(string(t.Side)),
				Type:   "market",
				Volume: t.Volume,
				//Price:  t.Price,
				Margin: string(t.Market),
			})

			timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
			payload := fmt.Sprintf("POST%s%s%s", endpoint, string(jsonBody), timestamp)
			sign := models.Sign(payload, models.Conf.Exchanges[string(t.Ex)].SecretKey, sha256.New, "hex")

			rq.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
			rq.Header.Add("Content-Type", "application/json")
			rq.Header.Add("X-COINEX-KEY", models.Conf.Exchanges[string(t.Ex)].ApiKey)
			rq.Header.Add("X-COINEX-SIGN", sign)
			rq.Header.Add("X-COINEX-TIMESTAMP", timestamp)

		},
		Params:   CoinexTradeParams{},
		Response: &TradeRes.CoinexTrade{},
	}
}
