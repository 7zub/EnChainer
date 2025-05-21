package OtherReq

import (
	"bytes"
	"crypto/sha256"
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/OtherRes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type CoinexTransferParams struct {
	From   string `url:"-" json:"from_account_type"`
	To     string `url:"-" json:"to_account_type"`
	Market string `url:"-" json:"market"`
	Ccy    string `url:"-" json:"ccy"`
	Amount string `url:"-" json:"amount"`
}

func (CoinexTransferParams) GetParams(task any) *models.Request {
	t := task.(models.TransferTask)
	endpoint := "/v2/assets/transfer"

	return &models.Request{
		Url:     "https://api.coinex.com" + endpoint,
		ReqType: models.ReqType.Transfer,
		SignWay: func(rq *http.Request) {

			jsonBody, _ := json.Marshal(CoinexTransferParams{
				From:   string(models.Spot),
				To:     string(models.Isolate),
				Market: t.Currency + t.Currency2,
				Ccy:    t.Currency2,
				Amount: fmt.Sprintf("%g", t.Amount),
			})

			fmt.Println(string(jsonBody))

			timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
			payload := fmt.Sprintf("POST%s%s%s", endpoint, string(jsonBody), timestamp)
			sign := models.Sign(payload, models.Conf.Exchanges[string(t.Ex)].SecretKey, sha256.New, "hex")

			rq.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
			rq.Header.Add("Content-Type", "application/json")
			rq.Header.Add("X-COINEX-KEY", models.Conf.Exchanges[string(t.Ex)].ApiKey)
			rq.Header.Add("X-COINEX-SIGN", sign)
			rq.Header.Add("X-COINEX-TIMESTAMP", timestamp)

		},
		Params:   CoinexTransferParams{},
		Response: &OtherRes.CoinexTransfer{},
	}
}
