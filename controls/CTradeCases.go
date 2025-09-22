package controls

import (
	"enchainer/models"
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
		slip = 0
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
	} else if opr.Ex == models.KUCOIN {
		decPrice = 3
		decVol = 1
	} else if opr.Ex == models.BINANCE && opr.Volume > 0.01 && opr.Volume < 1 {
		decPrice = 3
		decVol = 1
	}

	opr.Price = RoundSn(opr.Price, decPrice, mode, slip)

	if pend == false {
		opr.Volume = RoundSn(models.Const.Lot/opr.Price, decVol, "down", 0)
	}

	opr.CreateDate = time.Now()
}

func NeedContract(opr *models.OperationTask) {
	if !((opr.Ex == models.GATEIO || opr.Ex == models.HUOBI || opr.Ex == models.OKX) && opr.Market == models.Market.Futures) {
		return
	}

	var act, _ = CreateAction(*opr, models.ReqType.Contract)

	switch opr.Ex {
	case models.GATEIO, models.HUOBI:
		opr.Cct = PairInfo[opr.Ccy.Currency+"-"+string(opr.Ex)].Cct

	case models.OKX:
		opr.Cct = act.Any.(float64)
	}
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
