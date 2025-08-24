package TradeReq

import (
	"crypto/sha256"
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/TradeRes"
	"net/http"
	"strings"
	"time"
)

type MexcTradeParams struct {
	Ccy    string  `url:"symbol"`
	Side   string  `url:"side"`
	Type   string  `url:"type"`
	Volume float64 `url:"quantity"`
	Price  float64 `url:"price"`
	Live   string  `url:"timeInForce"`
	Time   int64   `url:"timestamp"`

	Contract string  `url:"contract"`
	IntSide  int     `url:"side"`
	Vol      float64 `url:"vol"`
	Leverage int     `url:"leverage"`
	IntType  int     `url:"type"`
	OpenType int     `url:"openType"`
}

func (MexcTradeParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)

	var url string
	var params MexcTradeParams
	switch t.Market {
	case models.Market.Spot:
		url = "https://api.mexc.com/api/v3/order"
		params = MexcTradeParams{
			Ccy:    t.Ccy.Currency + t.Ccy.Currency2,
			Side:   strings.ToUpper(string(t.Side)),
			Volume: t.Volume,
			Type:   "LIMIT",
		}
	case models.Market.Futures:
		url = "https://contract.mexc.com/api/v1/private/order/submit"
		params = MexcTradeParams{
			Contract: t.Ccy.Currency + "_" + t.Ccy.Currency2,
			IntSide:  3,
			Vol:      t.Volume,
			Leverage: 10,
			IntType:  1,
			OpenType: 2,
		}
	}

	return &models.Request{
		Url:     url,
		ReqType: models.ReqType.Trade,
		SignWay: func(rq *http.Request) {
			rq.Header.Add("X-MEXC-APIKEY", models.Conf.Exchanges[string(t.Ex)].ApiKey)
			q := rq.URL.Query()
			sign := models.Sign(rq.URL.Query().Encode(), models.Conf.Exchanges[string(t.Ex)].SecretKey, sha256.New, "hex")
			q.Add("signature", sign)
			rq.URL.RawQuery = q.Encode()
		},
		Params: MexcTradeParams{
			Ccy:    params.Ccy,
			Side:   params.Side,
			Type:   params.Type,
			Volume: params.Volume,
			Price:  t.Price,
			Live:   "GTC",
			Time:   time.Now().UnixMilli(),

			Contract: params.Contract,
			IntSide:  params.IntSide,
			Vol:      params.Vol,
			Leverage: params.Leverage,
			OpenType: params.OpenType,
		},
		Response: &TradeRes.MexcTrade{},
	}
}
