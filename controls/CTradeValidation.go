package controls

import (
	"enchainer/models"
	"fmt"
)

func TradeTaskValidation(task *models.TradeTask) {
	task.Stage = models.Validation

	if cnt >= models.Const.MaxTrade {
		task.Status = models.Stop
		task.Message += "Превышен лимит открытых тасок; "
	}

	if SearchOpenTask(task) != nil {
		task.Status = models.Stop
		task.Message += "Таска на пару уже существует; "
	}

	if task.Spread < models.Const.Spread {
		task.Status = models.Stop
		task.Message += "Низкий спред; "
	}

	if task.Buy.Price*task.Buy.Volume < models.Const.Lot {
		task.Status = models.Stop
		task.Message += fmt.Sprintf("Низкий объем на покупку: %g; ", task.Buy.Price*task.Buy.Volume)
	}

	if task.Sell.Price*task.Sell.Volume < models.Const.Lot {
		task.Status = models.Stop
		task.Message += fmt.Sprintf("Низкий объем на продажу: %g; ", task.Sell.Price*task.Sell.Volume)
	}

	if task.Sell.Ex == models.MEXC && task.Sell.Market == models.Market.Spot {
		task.Status = models.Stop
		task.Message += "На бирже MEXC нет маржинальной торговли; "
	}

	if (task.Buy.Ex == models.MEXC || task.Sell.Ex == models.MEXC) && task.Sell.Market == models.Market.Futures {
		task.Status = models.Stop
		task.Message += "Фьючерсная торговля на MEXC отключена; "
	}

	if task.Buy.Ex == models.OKX || task.Sell.Ex == models.OKX {
		task.Status = models.Stop
		task.Message += "Объем в OKX указан в контрактах; "
	}

	if task.Buy.Ex == models.GATEIO || task.Sell.Ex == models.GATEIO {
		task.Status = models.Stop
		task.Message += "Объем в GATEIO указан в контрактах; "
	}
}
