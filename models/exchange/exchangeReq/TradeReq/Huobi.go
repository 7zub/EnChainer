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
	Purpose   string  `url:"trade-purpose" json:"trade-purpose"`

	Contract string  `url:"contract_code" json:"contract_code"`
	Vol      float64 `url:"volume" json:"volume"`
	Dir      string  `url:"direction" json:"direction"`
	Offset   string  `url:"offset" json:"offset"`
	Lever    int     `url:"lever_rate" json:"lever_rate"`
	Mark     string  `url:"order_price_type" json:"order_price_type"`
}

func (HuobiTradeParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)

	var url, endpoint string
	var params HuobiTradeParams
	switch t.Market {
	case models.Market.Spot:
		url = "api.huobi.pro"
		endpoint = "/v1/order/orders/place"
		params = HuobiTradeParams{
			AccountId: "68842665", //68842529
			Ccy:       strings.ToLower(t.Currency + t.Currency2),
			Type:      string(t.Side) + "-limit",
			Volume:    t.Volume,
			Price:     t.Price,
			Source:    "super-margin-api",
			//Purpose:   "1",
		}
	case models.Market.Futures:
		url = "api.hbdm.com"
		endpoint = "/linear-swap-api/v1/swap_cross_order"
		params = HuobiTradeParams{
			Contract: strings.ToLower(t.Ccy.Currency + "-" + t.Ccy.Currency2),
			Vol:      t.Volume,
			Price:    t.Price,
			Dir:      string(t.Side),
			//Offset:   "open", в реж. хеджирования
			Lever: 10,
			Mark:  "limit",
		}
	}

	return &models.Request{
		Url:     "https://" + url + endpoint,
		ReqType: models.ReqType.Trade,
		SignWay: func(rq *http.Request) {
			jsonBody, _ := json.Marshal(params) //TODO возможно подтягивание лишних полей в json (нужен: ,omitempty)
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
		Response: &TradeRes.HuobiTrade{},
	}
}
