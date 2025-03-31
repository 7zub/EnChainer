package TradeRes

import (
	"enchainer/models"
)

type BybitTrade struct {
	OrderId string `json:"orderId"`
}

func (book BybitTrade) Mapper() any {
	if len(book.OrderId) > 0 {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
