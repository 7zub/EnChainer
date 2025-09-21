package controls

import (
	"enchainer/models"
	"fmt"
	"reflect"
	"sync"
)

var mu sync.Mutex
var maxTrade = 0
var activeTrade = 0

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
			Cct:       PairInfo[task.Ccy.Currency+"-"+string(task.Buy.Ex)].Cct,
		}
		oprSell := models.OperationTask{
			Ccy:       task.Ccy,
			Operation: task.Sell,
			Cct:       PairInfo[task.Ccy.Currency+"-"+string(task.Sell.Ex)].Cct,
		}

		PreparedOperation(&oprBuy, false)
		PreparedOperation(&oprSell, false)

		NeedContract(&oprBuy)
		NeedContract(&oprSell)

		ntBuy := NeedTransfer(&oprBuy, false)
		ntSell := NeedTransfer(&oprSell, false)

		if ntBuy.Status == models.OK && ntSell.Status == models.OK {
			var oSell, oBuy models.Result
			oSell, oprSell.ReqId = CreateAction(oprSell, models.ReqType.Trade)
			oBuy, oprBuy.ReqId = CreateAction(oprBuy, models.ReqType.Trade)

			if oSell.Status == models.OK && oBuy.Status == models.OK {
				task.Status = models.Pending
				mu.Lock()
				TaskTime(task.Ccy)
				//activeTrade += 1
				//if activeTrade == models.Const.MaxTrade {
				//	go TaskPause()
				//}
				mu.Unlock()
			} else {
				task.Status = models.Err
				task.Message = fmt.Sprintf("Ошибка открытия позиций: %s %s, %s %s", oprBuy.Side, oBuy.Status, oprSell.Side, oSell.Status)
			}
		} else {
			task.Status = models.Err
			task.Message = fmt.Sprintf("Ошибка трансфера: %s %s", ntBuy.Message, ntSell.Message)
		}

		task.OpTask = append(task.OpTask, oprSell, oprBuy)
		mu.Lock()
		maxTrade += 1
		mu.Unlock()
	}

	if task.Stage == models.Validation && task.Status == models.Stop {
		TradeTask.Delete(task.TaskId)
	} else {
		TradeTask.Store(task.TaskId, task)
	}
	ChanAny <- task
}

func CreateAction(act any, reqtype models.RqType) (models.Result, string) {
	ex := reflect.ValueOf(act).FieldByName("Ex").String()
	typ := GetTypeEx(models.Exchange(ex), reqtype)
	rr, _ := reflect.New(typ).Interface().(models.IParams)
	rq := rr.GetParams(act)
	rq.DescRequest(models.GenDescRequest())
	rq.SendRequest()
	ToLog(*rq)
	ChanAny <- rq
	res := rq.Response.Mapper().(models.Result)

	switch res.Status {
	case models.ERR:
		res.Message = fmt.Sprintf("%s %s не выполнен: %s", reqtype, rq.ReqId, rq.ResponseRaw)
	case models.OK:
		res.Message = fmt.Sprintf("%s %s выполнен: %s", reqtype, rq.ReqId, rq.ResponseRaw)
	}
	ToLog(res)
	return res, rq.ReqId
}
