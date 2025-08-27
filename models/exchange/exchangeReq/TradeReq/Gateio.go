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
	Ccy       string  `url:"-" json:"currency_pair,omitempty"`
	Side      string  `url:"-" json:"side"`
	Type      string  `url:"-" json:"type"`
	Volume    float64 `url:"-" json:"amount,omitempty"`
	Price     float64 `url:"-" json:"price"`
	Live      string  `url:"-" json:"time_in_force,omitempty"`
	Account   string  `url:"-" json:"account,omitempty"`
	Margin    string  `url:"-" json:"auto_borrow,omitempty"`
	AutoRepay string  `url:"-" json:"auto_repay,omitempty"`

	Contract string  `url:"-" json:"contract"`
	Size     float64 `url:"-" json:"size"`
	Live1    string  `url:"-" json:"tif"`
}

func (GateioTradeParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)

	var endpoint string
	var params GateioTradeParams
	switch t.Market {
	case models.Market.Spot:
		endpoint = "/api/v4/spot/orders"
		params = GateioTradeParams{
			Ccy:       t.Ccy.Currency + "_" + t.Ccy.Currency2,
			Side:      string(t.Side),
			Type:      "limit",
			Volume:    t.Volume,
			Live:      "gtc",
			Account:   "unified",
			Margin:    "true",
			AutoRepay: "true",
		}
	case models.Market.Futures:
		endpoint = "/api/v4/futures/" + strings.ToLower(t.Ccy.Currency2) + "/orders"
		params = GateioTradeParams{
			Contract: t.Ccy.Currency + "_" + t.Ccy.Currency2,
			Live1:    "gtc",
		}
		if t.Side == models.Sell {
			params.Size = -t.Volume
		} else {
			params.Size = t.Volume
		}
	}

	return &models.Request{
		Url:     "https://api.gateio.ws" + endpoint,
		ReqType: models.ReqType.Trade,
		SignWay: func(rq *http.Request) {
			jsonBody, _ := json.Marshal(GateioTradeParams{
				Ccy:       params.Ccy,
				Side:      params.Side,
				Type:      params.Type,
				Volume:    params.Volume,
				Price:     t.Price,
				Live:      params.Live,
				Account:   params.Account,
				Margin:    params.Margin,
				AutoRepay: params.AutoRepay,

				Contract: params.Contract,
				Size:     params.Size,
				Live1:    params.Live1,
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
