package controls

import (
	"enchainer/models"
	"fmt"
)

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

	if task.Sell.Ex == models.BINANCE && (task.Ccy.Currency == "KDA" || task.Ccy.Currency == "VANRY") {
		task.Status = models.Stop
		task.Message += "На бирже Binance не работает маржа; "
	}
}
