package TradeRes

import (
	"enchainer/models"
)

type MexcTrade struct {
	OrderId string `json:"orderId"`
}

func (book MexcTrade) Mapper() any {
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
