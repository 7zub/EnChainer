package controls

import (
	"enchainer/models"
	"fmt"
)

func PendingHandler(ccy models.Ccy, book []models.OrderBook) {
	if pendId := SearchPendTask(ccy); pendId != nil {
		task := LoadTask(*pendId)

		var i int
		var bid, ask float64
		var bkid []string
		for i = range book {
			if book[i].Exchange == task.Buy.Ex {
				bid = book[i].Bids[0].Price
				bkid = append(bkid, book[i].ReqId)
			}

			if book[i].Exchange == task.Sell.Ex {
				ask = book[i].Asks[0].Price
				bkid = append(bkid, book[i].ReqId)
			}
		}

		if bid <= 0 || ask <= 0 {
			ToLog(models.Result{
				Status: models.WAR,
				Message: fmt.Sprintf("Отсутствует книга: bids %f, ask %f для запросов %v, TaskId: %s, Ccy: %s",
					bid, ask, bkid, task.TaskId, task.Ccy.Currency)})
			return
		}

		profit := ((bid-task.Buy.Price)/task.Buy.Price + (task.Sell.Price-ask)/task.Sell.Price) * 100

		if profit < models.Const.MinProfit {
			ToLog(models.Result{
				Status: models.WAR,
				Message: fmt.Sprintf("Неприбыльный спред: %f, спред: %f, TaskId: %s, Ccy: %s",
					profit, task.Spread, task.TaskId, task.Ccy.Currency)})
			return
		}

		ToLog(models.Result{
			Status: models.WAR,
			Message: fmt.Sprintf("Возможность для прибыли: %f%%, спред: %f%%, TaskId: %s, Ccy: %s",
				profit, task.Spread, task.TaskId, task.Ccy.Currency),
		})

		opr1 := models.OperationTask{
			Ccy:       task.Ccy,
			Operation: task.OpTask[0].Operation,
			Cct:       task.OpTask[0].Cct,
		}

		opr2 := models.OperationTask{
			Ccy:       task.Ccy,
			Operation: task.OpTask[1].Operation,
			Cct:       task.OpTask[1].Cct,
		}

		opr1.Operation.Side = models.Buy
		opr1.Operation.Price = ask
		opr2.Operation.Side = models.Sell
		opr2.Operation.Price = bid

		PreparedOperation(&opr1, true)
		PreparedOperation(&opr2, true)

		var o1, o2 models.Result
		o1, opr1.ReqId = CreateAction(opr1, models.ReqType.Trade)
		o2, opr2.ReqId = CreateAction(opr2, models.ReqType.Trade)

		task.OpTask = append(task.OpTask, opr1, opr2)
		TradeTask.Store(task.TaskId, task)

		nt1 := NeedTransfer(&opr1, false)
		nt2 := NeedTransfer(&opr2, false)

		if o1.Status == models.OK && o2.Status == models.OK {
			if nt1.Status == models.OK && nt2.Status == models.OK {
				task.Status = models.Done
			} else {
				task.Status = models.Err
				task.Message = nt1.Message + nt2.Message
			}
		} else {
			task.Status = models.Err
			task.Message = fmt.Sprintf("Ошибка закрытия позиций: %s %s, %s %s", opr1.Side, o1.Status, opr2.Side, o2.Status)
		}

		TradeTask.Store(task.TaskId, task)
		ChanAny <- task
	}
}
