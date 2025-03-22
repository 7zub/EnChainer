package TradeRes

import (
	"enchainer/models"
)

type GateioTrade struct {
	OrderId string `json:"id"`
}

func (book GateioTrade) Mapper() any {
	if len(book.OrderId) > 1 {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
