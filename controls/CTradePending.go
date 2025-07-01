package controls

import (
	"enchainer/models"
	"fmt"
)

func PendingHandler(ccy models.Ccy, book []models.OrderBook) {
	if pendId := SearchPendTask(ccy); pendId != nil {
		task := LoadTask(*pendId)

		var b, s float64
		for i := range book {
			if book[i].Exchange == task.Buy.Ex {
				b = book[i].Asks[0].Price
			}

			if book[i].Exchange == task.Sell.Ex {
				s = book[i].Bids[0].Price
			}
		}

		if b <= 0 || s <= 0 {
			ToLog(models.Result{Status: models.WAR, Message: fmt.Sprintf("Отсутсвует книга: b %f, s %f", b, s)})
			return
		}

		spr := Round((s/b-1)*100, 4)

		if task.Spread-spr > models.Const.SpreadClose {
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

			fmt.Println("Проверка " + fmt.Sprintf("%d", len(task.OpTask)))

			if o1.Status == models.OK && o2.Status == models.OK {
				if NeedTransfer(&opr1, true).Status == models.OK && NeedTransfer(&opr2, true).Status == models.OK {
					task.Status = models.Done
				} else {
					task.Status = models.Err
					task.Message = "Ошибка трансфера в спот"
				}
			} else {
				task.Status = models.Err
				task.Message = "Ошибка закрытия операции: " + string(opr1.Side) + " " + string(o1.Status) + "; " + string(opr2.Side) + " " + string(o2.Status)
			}

			TradeTask.Store(task.TaskId, task)
			SaveDb(&task)
		} else {
			ToLog(models.Result{Status: models.WAR, Message: fmt.Sprintf("Маленькая разница с новым спредом: %f", task.Spread-spr)})
		}
	}
}
