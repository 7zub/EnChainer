package TradeRes

import (
	"enchainer/models"
)

type MexcTrade struct {
	OrderId int64 `json:"orderId"`
}

func (book MexcTrade) Mapper() any {
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
