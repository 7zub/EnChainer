package controls

import (
	"enchainer/models"
	"fmt"
	"reflect"
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

		task.OpTask = append(task.OpTask, oprSell, oprBuy)
		cnt += 1
	}

	if task.Stage == models.Validation && task.Status == models.Stop {
		TradeTask.Delete(task.TaskId)
	} else {
		TradeTask.Store(task.TaskId, task)
	}
	ChanAny <- task
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
	ChanAny <- rq
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
