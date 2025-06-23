package TradeRes

import (
	"enchainer/models"
)

type HuobiTrade struct {
	OrderId string `json:"data"`
}

func (book HuobiTrade) Mapper() any {
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
