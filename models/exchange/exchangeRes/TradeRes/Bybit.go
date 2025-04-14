package TradeRes

import (
	"enchainer/models"
)

type BybitTrade struct {
	Result struct {
		OrderId string `json:"orderId"`
	} `json:"result"`
}

func (book BybitTrade) Mapper() any {
	if len(book.Result.OrderId) > 0 {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
