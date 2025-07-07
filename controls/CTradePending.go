package controls

import (
	"enchainer/models"
	"fmt"
)

func PendingHandler(ccy models.Ccy, book []models.OrderBook) {
	if pendId := SearchPendTask(ccy); pendId != nil {
		task := LoadTask(*pendId)

		var i int
		var b, s float64
		for i = range book {
			if book[i].Exchange == task.Buy.Ex {
				b = book[i].Asks[0].Price
			}

			if book[i].Exchange == task.Sell.Ex {
				s = book[i].Bids[0].Price
			}
		}

		if b <= 0 || s <= 0 {
			ToLog(models.Result{Status: models.WAR, Message: fmt.Sprintf("Отсутсвует книга: bids %f, ask %f для %v", b, s, task)})
			return
		}

		spr := Round((s/b-1)*100, 4)

		if task.Spread-spr > models.Const.SpreadClose {
			ToLog(models.Result{Status: models.WAR, Message: fmt.Sprintf("Найден спрэд: %f, %v", task.Spread-spr, book[i])})

			opr1 := models.OperationTask{
				Ccy:       task.Ccy,
				Operation: task.OpTask[0].Operation,
			}

			opr2 := models.OperationTask{
				Ccy:       task.Ccy,
				Operation: task.OpTask[1].Operation,
			}

			opr1.Operation.Side = models.Sell
			opr1.Operation.Price = s
			opr2.Operation.Side = models.Buy
			opr2.Operation.Price = b

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
				task.Message = "Ошибка закрытия операции: " + string(opr1.Side) + " " + string(o1.Status) + "; " + string(opr2.Side) + " " + string(o2.Status)
			}

			TradeTask.Store(task.TaskId, task)
			SaveDb(&task)
		} else {
			ToLog(models.Result{Status: models.WAR, Message: fmt.Sprintf("Маленькая разница с новым спредом: %f, %v", task.Spread-spr, task)})
		}
	}
}
