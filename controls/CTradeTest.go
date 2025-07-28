package controls

import (
	"enchainer/models"
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
		Spread: 1,
		Buy: models.Operation{
			Ex:     models.COINEX,
			Price:  0.36,
			Volume: 60,
			Side:   models.Buy,
		},
		Sell: models.Operation{
			Ex:     models.GATEIO,
			Price:  0.3,
			Volume: 35.9,
			Side:   models.Sell,
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
		From: models.Spot,
		To:   models.Isolate,
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
