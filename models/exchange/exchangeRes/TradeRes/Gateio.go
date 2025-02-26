package TradeRes

import (
	"enchainer/models"
)

type GateioTrade struct {
	OrderId int64 `json:"id"`
}

func (book GateioTrade) Mapper() any {
	if book.OrderId > 0 {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
