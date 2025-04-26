package controls

import (
	"enchainer/models"
	"time"
)

func PreparedOperation(opr *models.OperationTask, pend bool) {
	var mode string
	var decPrice = 4
	var decVol = 3

	if opr.Side == models.Buy {
		mode = "up"
	} else {
		mode = "down"
	}

	if opr.Ex == models.BYBIT {
		if opr.Price > 0.2 && opr.Price < 1 {
			decPrice = 3
		}

		if 5.2/opr.Price > 3 && 5.2/opr.Price < 100 {
			decVol = 2
		}
	}

	opr.Price = RoundSn(opr.Price, decPrice, mode)

	if pend == false {
		opr.Volume = RoundSn(5.2/opr.Price, decVol, "down")
	}

	opr.CreateDate = time.Now()
}
