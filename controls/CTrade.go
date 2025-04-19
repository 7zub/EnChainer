package controls

import (
	"enchainer/models"
	"fmt"
	"reflect"
	"time"
)

var cnt = 0

func TradeTaskHandler(task *models.TradeTask) {
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
		oprBuy.Price = RoundSn(oprBuy.Price, 4)
		oprSell.Price = RoundSn(oprSell.Price, 4)
		oprBuy.Volume = RoundSn(5.2/oprBuy.Price, 3)
		oprSell.Volume = RoundSn(5.2/oprSell.Price, 3)

		if oprBuy.Ex == models.BYBIT {
			if oprBuy.Price < 1 {
				oprBuy.Price = RoundSn(oprBuy.Price, 3)
			}

			if oprBuy.Volume > 10 && oprBuy.Volume < 100 {
				oprBuy.Volume = RoundSn(oprBuy.Volume, 2)
			}
		}

		if oprSell.Ex == models.BYBIT {
			if oprSell.Price < 1 {
				oprSell.Price = RoundSn(oprSell.Price, 3)
			}

			if oprSell.Volume > 10 && oprSell.Volume < 100 {
				oprSell.Volume = RoundSn(oprSell.Volume, 2)
			}
		}

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

	TradeTask.Store(task.TaskId, task)
	SaveDb(&task)
}

func TradeTaskValidation(task *models.TradeTask) {
	task.Stage = models.Validation

	if cnt >= 1 {
		task.Status = models.Stop
		task.Message += "Превышен лимит открытых тасок; "
	}

	if SearchOpenTask(task) != nil {
		task.Status = models.Stop
		task.Message += "Таска на пару уже существует; "
	}

	if task.Spread < 0.3 {
		task.Status = models.Stop
		task.Message += "Низкий спред; "
	}

	if task.Buy.Price*task.Buy.Volume < 5 {
		task.Status = models.Stop
		task.Message += fmt.Sprintf("Низкий объем на покупку: %g; ", task.Buy.Price*task.Buy.Volume)
	}

	if task.Sell.Price*task.Sell.Volume < 5 {
		task.Status = models.Stop
		task.Message += fmt.Sprintf("Низкий объем на продажу: %g; ", task.Sell.Price*task.Sell.Volume)
	}

	if task.Sell.Ex == models.MEXC {
		task.Status = models.Stop
		task.Message += "На бирже MEXC нет маржинальной торговли; "
	}
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
			Ex:     models.GATEIO,
			Price:  70,
			Volume: 0.09,
			Side:   models.Buy,
		},
		Sell: models.Operation{
			Ex:     models.BINANCE,
			Price:  200,
			Volume: 0.03,
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
