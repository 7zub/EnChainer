package TradeRes

import (
	"enchainer/models"
)

type GateioTrade struct {
	OrderId  string `json:"id"`
	Contract string `json:"contract"`
}

func (book GateioTrade) Mapper() any {
	if len(book.OrderId) > 1 || len(book.Contract) > 1 {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
