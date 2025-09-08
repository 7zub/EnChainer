package controls

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/ContractRes"
	"strconv"
	"time"
)

func PreparedOperation(opr *models.OperationTask, pend bool) {
	var mode string
	var slip float64
	var decPrice = models.Const.DecimalPrice
	var decVol = models.Const.DecimalVolume

	if opr.Side == models.Buy {
		mode = "up"
	} else {
		mode = "down"
	}

	if opr.Ex == models.KUCOIN {
		slip = models.Const.Slip * 5
	} else if opr.Deep >= 3 {
		slip = models.Const.Slip * 3
	} else if opr.Deep > 1 {
		slip = models.Const.Slip * 2
	} else {
		slip = models.Const.Slip
	}

	if opr.Ex == models.BYBIT && opr.Market == models.Market.Spot {
		if opr.Price > 0.2 && opr.Price < 1 {
			decPrice = 3
		}

		if models.Const.Lot/opr.Price > 3 && models.Const.Lot/opr.Price < 100 {
			decVol = 2
		}
	}

	opr.Price = RoundSn(opr.Price, decPrice, mode, slip)

	if pend == false {
		opr.Volume = RoundSn(models.Const.Lot/opr.Price, decVol, "down", 0)
	}

	opr.CreateDate = time.Now()
}

func NeedTransfer(opr *models.OperationTask, isl bool) models.Result {
	if opr.Ex != models.COINEX || opr.Market == models.Market.Futures {
		return models.Result{Status: models.OK}
	}

	from, to := models.Market.Spot, models.Market.Isolate
	if isl {
		from, to = to, from
	}

	var tr = models.TransferTask{
		Ex:         opr.Ex,
		From:       from,
		To:         to,
		Ccy:        opr.Ccy,
		Amount:     models.Const.Lot,
		CreateDate: time.Now(),
	}

	trf, _ := CreateAction(tr, models.ReqType.Transfer)
	return trf
}

func NeedContract(opr *models.OperationTask) models.Result {
	if !((opr.Ex == models.GATEIO || opr.Ex == models.OKX || opr.Ex == models.HUOBI) && opr.Market == models.Market.Futures) {
		return models.Result{Status: models.OK}
	}

	var act, _ = CreateAction(*opr, models.ReqType.Contract)

	switch opr.Ex {
	case models.GATEIO:
		for _, c := range act.Any.(ContractRes.GateioContract) {
			if c.Ccy == opr.Ccy.Currency+"_"+opr.Ccy.Currency2 {
				opr.Cct, _ = strconv.ParseFloat(c.Cct, 64)
				if opr.Cct <= 0 {
					return models.Result{Status: models.ERR, Message: "Ошибка получения контракта"}
				}
				return act
			}
		}

	case models.OKX, models.HUOBI:
		opr.Cct = act.Any.(float64)
	}

	return models.Result{Status: models.ERR, Message: "Не найдена валюта"}
}
