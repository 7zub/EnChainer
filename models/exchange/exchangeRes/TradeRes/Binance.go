package TradeRes

import (
	"enchainer/models"
)

type BinanceTrade struct {
	OrderId int64 `json:"orderId"`
}

func (book BinanceTrade) Mapper() any {
	if book.OrderId > 1 {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
