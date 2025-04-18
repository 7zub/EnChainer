package controls

import (
	"enchainer/models"
	"fmt"
)

func PendingHandler(ccy models.Ccy, book models.OrderBook) {
	if pendId, i := SearchOperation(ccy, book.Exchange); pendId != nil {
		task := LoadTask(*pendId)

		opr := models.OperationTask{
			Ccy:       task.Ccy,
			Operation: task.OpTask[i].Operation,
		}

		switch i {
		case 0:
			opr.Operation.Side = models.Sell
			opr.Operation.Price = RoundSn(book.Bids[0].Price, 4)
		case 1:
			opr.Operation.Side = models.Buy
			opr.Operation.Price = RoundSn(book.Asks[0].Price, 4)
		}
		var o models.Result
		o, opr.ReqId = CreateOrder(opr)
		task.Mu.Lock()
		task.OpTask = append(task.OpTask, opr)
		task.Mu.Unlock()
		TradeTask.Store(task.TaskId, task)
		l := len(task.OpTask)

		fmt.Println("Проверка " + fmt.Sprintf("%d", len(task.OpTask)))

		if o.Status == models.OK {
			switch l {
			case 3:
				task.Status = models.Progress
			case 4:
				task.Status = models.Done
			default:
				task.Status = models.Err
				task.Message = fmt.Sprintf("Некорректное количество операций: %d", len(task.OpTask))
			}
		} else {
			task.Status = models.Err
			task.Message = "Ошибка завершения сделки: " + string(opr.Operation.Side) + " " + string(o.Status)
		}

		TradeTask.Store(task.TaskId, task)
		SaveDb(&task)
	}
}
