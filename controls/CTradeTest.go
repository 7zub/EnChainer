package controls

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeReq/LeverReq"
	"enchainer/models/exchange/exchangeReq/OtherReq"
	"fmt"
	"time"
)

func Trade() {
	task := models.TradeTask{
		TaskId: "10",
		Ccy: models.Ccy{
			Currency:  "FLOW",
			Currency2: "USDT",
		},
		Spread: 0.9,
		Buy: models.Operation{
			Ex:     models.COINEX,
			Price:  0.11,
			Volume: 300,
			Side:   models.Buy,
			Market: models.Market.Futures,
		},
		Sell: models.Operation{
			Ex:     models.OKX,
			Price:  0.6,
			Volume: 300,
			Side:   models.Sell,
			Market: models.Market.Futures,
		},
		CreateDate: time.Time{},
		Stage:      models.Creation,
		Status:     models.Done,
	}

	TradeTask.Store(task.TaskId, &task)
	SaveDb(&task)

	TradeTaskHandler(&task)
	fmt.Println(task.Stage)
}

func Trans() {
	var tr = models.TransferTask{
		Ex:   models.COINEX,
		From: models.Market.Spot,
		To:   models.Market.Isolate,
		Ccy: models.Ccy{
			Currency:  "SOL",
			Currency2: "USDT",
		},
		Amount:     5,
		CreateDate: time.Now(),
	}

	rr := OtherReq.CoinexTransferParams{}
	rq := rr.GetParams(tr)
	rq.DescRequest(models.GenDescRequest())
	rq.SendRequest()
	ToLog(*rq)
	go SaveDb(rq)
	SaveDb(&tr)
}

func Lever() {
	var tr = models.OperationTask{
		Ccy: models.Ccy{
			Currency:  "FLOW",
			Currency2: "USDT",
		},
		Operation: models.Operation{
			Ex:     models.BINANCE,
			Market: models.Market.Futures,
		},
	}

	rr := LeverReq.BinanceLeverageParams{
		Leverage: 20,
	}
	rq := rr.GetParams(tr)
	rq.DescRequest(models.GenDescRequest())
	rq.SendRequest()
	ToLog(*rq)
	ChanAny <- rq
}

func Contract() {
	oo := models.Operation{Ex: models.GATEIO, Market: models.Market.Futures}
	ooo := models.OperationTask{
		Ccy:       models.Ccy{Currency: "LUNC", Currency2: "USDT"},
		Operation: oo,
	}
	NeedContract(&ooo)
}
