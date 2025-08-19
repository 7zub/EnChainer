package LeverReq

import (
	"crypto/sha256"
	"enchainer/models"
	"net/http"
	"time"
)

type BinanceLeverageParams struct {
	Symbol    string `url:"symbol"`
	Leverage  int    `url:"leverage"`
	Timestamp int64  `url:"timestamp"`
}

func (BinanceLeverageParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)

	url := ""
	if t.Market == models.Market.Features {
		url = "https://fapi.binance.com/fapi/v1/leverage"
	}

	return &models.Request{
		Url:     url,
		ReqType: models.ReqType.Trade,
		SignWay: func(rq *http.Request) {
			rq.Header.Add("X-MBX-APIKEY", models.Conf.Exchanges[string(t.Ex)].ApiKey)
			q := rq.URL.Query()
			sign := models.Sign(rq.URL.Query().Encode(), models.Conf.Exchanges[string(t.Ex)].SecretKey, sha256.New, "hex")
			q.Add("signature", sign)
			rq.URL.RawQuery = q.Encode()
		},
		Params: BinanceLeverageParams{
			Symbol:    t.Currency + t.Currency2,
			Leverage:  50,
			Timestamp: time.Now().UnixMilli(),
		},
	}
}
