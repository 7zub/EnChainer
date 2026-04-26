package controls

import (
	"enchainer/controls/load"
	"enchainer/models"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

var maxTrade atomic.Int32
var activeTrade atomic.Int32

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

		NeedContract(&oprBuy)
		NeedContract(&oprSell)

		ntBuy := NeedTransfer(&oprBuy, false)
		ntSell := NeedTransfer(&oprSell, false)

		activeTrade.Add(1)

		if ntBuy.Status == models.OK && ntSell.Status == models.OK {
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				oprSell.Result, oprSell.ReqId = CreateAction(oprSell, models.ReqType.Trade)
			}()

			go func() {
				defer wg.Done()
				oprBuy.Result, oprBuy.ReqId = CreateAction(oprBuy, models.ReqType.Trade)
			}()
			wg.Wait()

			if oprSell.Status == models.OK && oprBuy.Status == models.OK {
				task.Status = models.Pending
				TaskTime(task.Ccy, 2)
			} else {
				task.Status = models.Err
				if op := TradeCancel(oprSell, oprBuy); op != nil {
					task.OpTask = append(task.OpTask, *op)
					task.Status = models.Cancel
				}

				task.Message = fmt.Sprintf("Ошибка открытия позиций: %s %s, %s %s", oprBuy.Side, oprBuy.Status, oprSell.Side, oprSell.Status)
			}
		} else {
			task.Status = models.Err
			task.Message = fmt.Sprintf("Ошибка трансфера: %s %s", ntBuy.Message, ntSell.Message)
		}

		task.OpTask = append(task.OpTask, oprSell, oprBuy)
		maxTrade.Add(1)
	}

	if task.Stage == models.Validation && task.Status == models.Stop {
		TradeTask.Delete(task.TaskId)
	} else {
		TradeTask.Store(task.TaskId, task)
	}
	ChanAny <- task
}

func TradeCancel(oprSell, oprBuy models.OperationTask) *models.OperationTask {
	var cancel models.OperationTask
	if oprSell.Status != models.OK && oprBuy.Status == models.OK {
		cancel = oprBuy
	} else if oprSell.Status == models.OK && oprBuy.Status != models.OK {
		cancel = oprSell
	} else {
		return nil
	}
	cancel.Side.Opposite()
	cancel.Result, cancel.ReqId = CreateAction(cancel, models.ReqType.Trade)
	return &cancel
}

func CreateAction(act any, reqType models.RqType) (models.Result, string) {
	ex := reflect.ValueOf(act).FieldByName("Ex").String()
	typ := GetTypeEx(models.Exchange(ex), reqType)
	rr, _ := reflect.New(typ).Interface().(models.IParams)
	rq := rr.GetParams(act)
	rq.DescRequest(models.GenDescRequest())
	rq.SendRequest()
	load.ToLog(rq.Log)
	ChanAny <- rq
	res := rq.Response.Mapper().(models.Result)

	switch res.Status {
	case models.ERR:
		res.Message = fmt.Sprintf("%s %s не выполнен: %s", reqType, rq.ReqId, rq.ResponseRaw)
	case models.OK:
		res.Message = fmt.Sprintf("%s %s выполнен: %s", reqType, rq.ReqId, rq.ResponseRaw)
	}
	load.ToLog(res)
	return res, rq.ReqId
}
