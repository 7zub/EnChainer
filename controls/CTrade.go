package controls

import (
	"enchainer/models"
	"fmt"
	"reflect"
	"time"
)

var cnt = 0

func TradeTaskHandler(task *models.TradeTask) {
	task.Mu.Lock()
	defer task.Mu.Unlock()

	if task.Stage == models.Creation && task.Status == models.Done {
		TradeTaskValidation(task)
	}

	if task.Stage == models.Validation && task.Status == models.Done {
		task.Stage = models.Trade

		oprBuy := models.OperationTask{
			Ccy:       task.Ccy,
			Operation: task.Buy,
		}
		oprSell := models.OperationTask{
			Ccy:       task.Ccy,
			Operation: task.Sell,
		}

		PreparedOperation(&oprBuy, false)
		PreparedOperation(&oprSell, false)

		var oSell, oBuy models.Result
		if oSell, oprSell.ReqId = CreateOrder(oprSell); oSell.Status == models.OK {
			if oBuy, oprBuy.ReqId = CreateOrder(oprBuy); oBuy.Status == models.OK {
				task.Status = models.Pending
			} else {
				task.Status = models.Err
				task.Message = "Ошибка операции: " + string(oprBuy.Side) + " " + string(oBuy.Status) + "; " + string(oprSell.Side) + " " + string(oSell.Status)
			}
		} else {
			task.Status = models.Err
			task.Message = "Ошибка операции: " + string(oprSell.Side)
		}

		task.OpTask = append(task.OpTask, oprBuy, oprSell)
		cnt += 1
	}

	if task.Stage == models.Validation && task.Status == models.Stop {
		TradeTask.Delete(task.TaskId)
	} else {
		TradeTask.Store(task.TaskId, task)
	}
	SaveDb(&task)
}

func CreateOrder(opr models.OperationTask) (models.Result, string) {
	typ := GetTypeEx(opr.Ex, "Trade")
	rr, _ := reflect.New(typ).Interface().(models.IParams)
	rq := rr.GetParams(opr)
	rq.DescRequest(models.GenDescRequest())
	rq.SendRequest()
	ToLog(*rq)
	go SaveDb(rq)
	res := rq.Response.Mapper().(models.Result)

	switch res.Status {
	case models.ERR:
		res.Message = "Операция " + rq.ReqId + " не выполнена: " + rq.ResponseRaw
	case models.OK:
		res.Message = "Операция " + rq.ReqId + " выполнена"
	}
	ToLog(res)
	return res, rq.ReqId
}

func Trade() {
	task := models.TradeTask{
		TaskId: "10",
		Ccy: models.Ccy{
			Currency:  "SOL",
			Currency2: "USDT",
		},
		Spread: 1,
		Buy: models.Operation{
			Ex:     models.BYBIT,
			Price:  70,
			Volume: 0.09,
			Side:   models.Buy,
		},
		Sell: models.Operation{
			Ex:     models.OKX,
			Price:  200,
			Volume: 0.05,
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
	//fmt.Println(TradeTask[0].Stage)
}
