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
			Ex:     models.GATEIO,
			Price:  0.25,
			Volume: 100,
			Side:   models.Buy,
			Market: models.Market.Features,
		},
		Sell: models.Operation{
			Ex:     models.BINANCE,
			Price:  0.9,
			Volume: 100,
			Side:   models.Sell,
			Market: models.Market.Features,
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
			Currency:  "SOL",
			Currency2: "USDT",
		},
		Operation: models.Operation{
			Ex:     models.BINANCE,
			Market: models.Market.Features,
		},
	}

	rr := LeverReq.BinanceLeverageParams{
		Leverage: 50,
	}
	rq := rr.GetParams(tr)
	rq.DescRequest(models.GenDescRequest())
	rq.SendRequest()
	ToLog(*rq)
	ChanAny <- rq
}
