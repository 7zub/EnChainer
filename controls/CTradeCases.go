package controls

import (
	"enchainer/models"
	"time"
)

func PreparedOperation(opr *models.OperationTask, pend bool) {
	var mode string
	var decPrice = models.Const.DecimalPrice
	var decVol = models.Const.DecimalVolume

	if opr.Side == models.Buy {
		mode = "up"
	} else {
		mode = "down"
	}

	if opr.Ex == models.BYBIT {
		if opr.Price > 0.2 && opr.Price < 1 {
			decPrice = 3
		}

		if models.Const.Lot/opr.Price > 3 && models.Const.Lot/opr.Price < 100 {
			decVol = 2
		}
	}

	opr.Price = RoundSn(opr.Price, decPrice, mode)

	if pend == false {
		opr.Volume = RoundSn(models.Const.Lot/opr.Price, decVol, "down")
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
