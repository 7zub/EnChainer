package TradeReq

import (
	"bytes"
	"crypto/sha512"
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/TradeRes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GateioTradeParams struct {
	Ccy       string  `url:"-" json:"currency_pair"`
	Side      string  `url:"-" json:"side"`
	Type      string  `url:"-" json:"type"`
	Volume    float64 `url:"-" json:"amount"`
	Price     float64 `url:"-" json:"price"`
	Live      string  `url:"-" json:"time_in_force"`
	Account   string  `url:"-" json:"account"`
	Margin    string  `url:"-" json:"auto_borrow"`
	AutoRepay string  `url:"-" json:"auto_repay"`
}

func (GateioTradeParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)
	endpoint := "/api/v4/spot/orders"

	return &models.Request{
		Url:     "https://api.gateio.ws" + endpoint,
		ReqType: models.ReqType.Trade,
		SignWay: func(rq *http.Request) {
			jsonBody, _ := json.Marshal(GateioTradeParams{
				Ccy:       t.Currency + "_" + t.Currency2,
				Side:      strings.ToLower(string(t.Side)),
				Type:      "limit",
				Volume:    t.Volume,
				Price:     t.Price,
				Live:      "gtc",
				Account:   "unified",
				Margin:    "true",
				AutoRepay: "true",
			})

			hash := sha512.Sum512(jsonBody)
			encBody := hex.EncodeToString(hash[:])
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			payload := fmt.Sprintf("POST\n%s\n\n%s\n%s", endpoint, encBody, timestamp)
			sign := models.Sign(payload, models.Conf.Exchanges[string(t.Ex)].SecretKey, sha512.New, "hex")

			rq.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
			rq.Header.Add("Accept", "application/json")
			rq.Header.Add("Content-Type", "application/json")
			rq.Header.Add("KEY", models.Conf.Exchanges[string(t.Ex)].ApiKey)
			rq.Header.Add("SIGN", sign)
			rq.Header.Add("Timestamp", timestamp)
		},
		Params:   GateioTradeParams{},
		Response: &TradeRes.GateioTrade{},
	}
}
