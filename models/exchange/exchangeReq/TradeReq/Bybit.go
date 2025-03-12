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

type BybitTradeParams struct {
	Ccy     string `url:"-" json:"symbol"`
	Side    string `url:"-" json:"side"`
	Type    string `url:"-" json:"orderType"`
	Volume  string `url:"-" json:"qty"`
	Price   string `url:"-" json:"price"`
	Live    string `url:"-" json:"timeInForce"`
	Account string `url:"-" json:"category"`
	Margin  int    `url:"-" json:"isLeverage"`
}

func (BybitTradeParams) GetParams(task any) *models.Request {
	t := task.(models.TradeTask)

	return &models.Request{
		Url:     "https://api.bybit.com/v5/order/create",
		ReqType: "Trade",
		SignWay: func(rq *http.Request) {

			jsonBody, _ := json.Marshal(BybitTradeParams{
				Ccy:     t.Ccy.Currency + t.Ccy.Currency2,
				Side:    strings.ToUpper(string(t.Stage[0])) + string(t.Stage[1:]),
				Type:    "Limit",
				Volume:  fmt.Sprintf("%g", t.Buy.Volume),
				Price:   fmt.Sprintf("%g", t.Buy.Price),
				Live:    "GTC",
				Account: "spot",
				Margin:  1,
			})

			timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
			sign := models.Sign(timestamp+models.Conf.Exchanges[t.Buy.Ex].ApiKey+string(jsonBody[:]), models.Conf.Exchanges[t.Buy.Ex].SecretKey, sha256.New)

			rq.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
			rq.Header.Set("Content-Type", "application/json")
			rq.Header.Set("X-BAPI-API-KEY", models.Conf.Exchanges[t.Buy.Ex].ApiKey)
			rq.Header.Set("X-BAPI-SIGN", sign)
			rq.Header.Set("X-BAPI-TIMESTAMP", timestamp)
		},
		Params:   BybitTradeParams{},
		Response: &TradeRes.BybitTrade{},
	}
}
