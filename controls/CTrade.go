package controls

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeReq/TradeReq"
	"time"
)

func createOrder1(task models.TradeTask) {
	//fmt.Println("req.URL")

	var rr TradeReq.BinanceTradeParams
	rq := rr.GetParams(task)
	rq.DescRequest(models.GenDescRequest())
	rq.SendRequest()
	ToLog(*rq)
	go SaveReqDb(rq)
}

func Trade() {
	task := models.TradeTask{
		TaskId: 10,
		Ccy: models.Ccy{
			Currency:  "SOL",
			Currency2: "USDT",
		},
		Spread: 0,
		Buy: models.Operation{
			Ex:     models.BINANCE,
			Price:  76,
			Volume: 0.1,
		},
		Sell: models.Operation{
			Ex:     "BY",
			Price:  0,
			Volume: 0,
		},
		CreateDate: time.Time{},
		Stage:      "buy",
		Status:     models.New,
	}

	createOrder1(task)
}
