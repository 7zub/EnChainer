package controls

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeReq/OtherReq"
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

		ntBuy := NeedTransfer(&oprBuy, false)
		ntSell := NeedTransfer(&oprSell, false)

		if ntBuy.Status == models.OK && ntSell.Status == models.OK {
			var oSell, oBuy models.Result
			if oSell, oprSell.ReqId = CreateAction(oprSell, models.ReqType.Trade); oSell.Status == models.OK {
				if oBuy, oprBuy.ReqId = CreateAction(oprBuy, models.ReqType.Trade); oBuy.Status == models.OK {
					task.Status = models.Pending
				} else {
					task.Status = models.Err
					task.Message = "Ошибка операции: " + string(oprBuy.Side) + " " + string(oBuy.Status) + "; " + string(oprSell.Side) + " " + string(oSell.Status)
				}
			} else {
				task.Status = models.Err
				task.Message = "Ошибка операции: " + string(oprSell.Side)
			}
		} else {
			task.Status = models.Err
			task.Message = ntBuy.Message + ntSell.Message
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

func CreateAction(act any, reqtype models.RqType) (models.Result, string) {
	v := reflect.ValueOf(act)
	ex := v.FieldByName("Ex")
	typ := GetTypeEx(models.Exchange(ex.String()), string(reqtype))
	rr, _ := reflect.New(typ).Interface().(models.IParams)
	rq := rr.GetParams(act)
	rq.DescRequest(models.GenDescRequest())
	rq.SendRequest()
	ToLog(*rq)
	go SaveDb(rq)
	res := rq.Response.Mapper().(models.Result)

	switch res.Status {
	case models.ERR:
		res.Message = fmt.Sprintf("%s %s не выполнен: %s", reqtype, rq.ReqId, rq.ResponseRaw)
	case models.OK:
		res.Message = fmt.Sprintf("%s %s успешна: %s", reqtype, rq.ReqId, rq.ResponseRaw)
	}
	ToLog(res)
	return res, rq.ReqId
}

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
	//fmt.Println(TradeTask[0].Stage)
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
