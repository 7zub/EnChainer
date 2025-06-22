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

type HuobiTradeParams struct {
	Ccy       string  `url:"symbol" json:"symbol"`
	Type      string  `url:"type" json:"type"`
	Volume    float64 `url:"amount" json:"amount"`
	Price     float64 `url:"price" json:"price"`
	AccountId string  `url:"account-id" json:"account-id"`
	Source    string  `url:"source" json:"source"`
}

func (HuobiTradeParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)
	url := "api.huobi.pro"
	endpoint := "/v1/order/orders/place"

	params := HuobiTradeParams{
		AccountId: "68842665", //68842529
		Ccy:       strings.ToLower(t.Currency + t.Currency2),
		Type:      strings.ToLower(string(t.Side)) + "-limit",
		Volume:    t.Volume,
		Price:     t.Price,
		Source:    "super-margin-api",
	}

	return &models.Request{
		Url:     "https://" + url + endpoint,
		ReqType: models.ReqType.Trade,
		SignWay: func(rq *http.Request) {
			jsonBody, _ := json.Marshal(params)
			timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")
			q := rq.URL.Query()
			q.Set("AccessKeyId", models.Conf.Exchanges[string(t.Ex)].ApiKey)
			q.Set("SignatureMethod", "HmacSHA256")
			q.Set("SignatureVersion", "2")
			q.Set("Timestamp", timestamp)
			payload := fmt.Sprintf("POST\n%s\n%s\n%s", url, endpoint, q.Encode())
			sign := models.Sign(payload, models.Conf.Exchanges[string(t.Ex)].SecretKey, sha256.New, "base64")
			q.Set("Signature", sign)
			rq.URL.RawQuery = q.Encode()

			rq.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
			rq.Header.Add("Content-Type", "application/json")

		},
		Params:   params,
		Response: &TradeRes.GateioTrade{},
	}
}
