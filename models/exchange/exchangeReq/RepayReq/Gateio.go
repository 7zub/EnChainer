package RepayReq

import (
	"bytes"
	"crypto/sha512"
	"enchainer/models"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type GateioRepayParams struct {
	Ccy       string  `url:"-" json:"currency"`
	Type      string  `url:"-" json:"type"`
	Volume    float64 `url:"-" json:"amount"`
	CcyPair   string  `url:"-" json:"currency_pair"`
	RepaidAll bool    `url:"-" json:"repaid_all"`
}

func (GateioRepayParams) GetParams(task any) *models.Request {
	t := task.(models.TradeTask)
	endpoint := "/api/v4/margin/uni/loans"

	return &models.Request{
		Url:     "https://api.gateio.ws" + endpoint,
		ReqType: "Trade",
		SignWay: func(rq *http.Request) {

			jsonBody, _ := json.Marshal(GateioRepayParams{
				CcyPair:   t.Ccy.Currency + "_" + t.Ccy.Currency2,
				Ccy:       t.Ccy.Currency,
				Type:      "repay",
				Volume:    t.Buy.Volume,
				RepaidAll: true,
			})

			fmt.Println(string(jsonBody))

			hash := sha512.Sum512(jsonBody)
			encBody := hex.EncodeToString(hash[:])
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			payload := fmt.Sprintf("POST\n%s\n\n%s\n%s", endpoint, encBody, timestamp)
			sign := models.Sign(payload, models.Conf.Exchanges[string(t.Buy.Ex)].SecretKey, sha512.New, "hex")

			rq.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
			rq.Header.Add("Accept", "application/json")
			rq.Header.Add("Content-Type", "application/json")
			rq.Header.Add("KEY", models.Conf.Exchanges[string(t.Buy.Ex)].ApiKey)
			rq.Header.Add("SIGN", sign)
			rq.Header.Add("Timestamp", timestamp)
		},
		Params: GateioRepayParams{},
	}
}
