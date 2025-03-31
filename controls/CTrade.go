package controls

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeReq/TradeReq"
	"fmt"
	"time"
)

var cnt int = 0

func TradeTaskHandler(task *models.TradeTask) {
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

		oprBuy.Operation.Volume = Round(3.2 / oprBuy.Operation.Price)
		oprSell.Operation.Volume = Round(3.2 / oprSell.Operation.Price)

		oBuy := CreateOrder(oprBuy).Status
		oSell := CreateOrder(oprSell).Status

		if oBuy == models.OK && oSell == models.OK {
			task.Status = models.Done
		} else {
			task.Status = models.Err
			task.Message = "Ошибка операции: покупка " + string(oBuy) + ", продажа " + string(oSell)
		}

		cnt += 1
	}

	SaveTradeTaskDb(task)
}

func TradeTaskValidation(task *models.TradeTask) {
	task.Stage = models.Validation

	if cnt > 3 {
		task.Status = models.Stop
		task.Message += "Превышен лимит открытых тасок; "
	}

	if SearchOpenTask(*task) > -1 {
		task.Status = models.Stop
		task.Message += "Таска на пару уже существует; "
	}

	if task.Spread < 0.5 {
		task.Status = models.Stop
		task.Message += "Низкий спред; "
	}

	if task.Buy.Price*task.Buy.Volume < 3 {
		task.Status = models.Stop
		task.Message += fmt.Sprintf("Низкий объем на покупку: %g; ", task.Buy.Price*task.Buy.Volume)
	}

	if task.Sell.Price*task.Sell.Volume < 3 {
		task.Status = models.Stop
		task.Message += fmt.Sprintf("Низкий объем на продажу: %g; ", task.Sell.Price*task.Sell.Volume)
	}

	if task.Sell.Ex == models.MEXC {
		task.Status = models.Stop
		task.Message += "На бирже MEXC нет маржинальной торговли; "
	}
}

func CreateOrder(opr models.OperationTask) models.Result {
	rr := getTradeReq(opr.Ex)
	rq := rr.GetParams(opr)
	rq.DescRequest(models.GenDescRequest())
	rq.SendRequest()
	ToLog(*rq)
	go SaveReqDb(rq)
	return rq.Response.Mapper().(models.Result)
}

func getTradeReq(ex models.Exchange) models.IParams {
	switch ex {
	case models.BINANCE:
		return TradeReq.BinanceTradeParams{}
	case models.GATEIO:
		return TradeReq.GateioTradeParams{}
	//case models.HUOBI:
	//	return TradeReq.HuobiTradeParams{}
	//case models.OKX:
	//	return TradeReq.OkxTradeParams{}
	case models.MEXC:
		return TradeReq.MexcTradeParams{}
	case models.BYBIT:
		return TradeReq.BybitTradeParams{}
	//case models.KUCOIN:
	//	return TradeReq.KucoinTradeParams{}
	default:
		return nil
	}
}

func Trade() {
	task := models.TradeTask{
		TaskId: 10,
		Ccy: models.Ccy{
			Currency:  "SOL",
			Currency2: "USDT",
		},
		Spread: 1,
		Buy: models.Operation{
			Ex:     models.MEXC,
			Price:  90,
			Volume: 0.02,
			Side:   models.Buy,
		},
		Sell: models.Operation{
			Ex:     models.GATEIO,
			Price:  200,
			Volume: 0.03,
			Side:   models.Sell,
		},
		CreateDate: time.Time{},
		Stage:      models.Creation,
		Status:     models.Done,
	}

	TradeTask = append(TradeTask, task)

	TradeTaskHandler(&TradeTask[0])
	fmt.Println(task.Stage)
	//fmt.Println(TradeTask[0].Stage)
}
