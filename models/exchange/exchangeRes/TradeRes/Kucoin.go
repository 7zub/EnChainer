package TradeRes

import (
	"enchainer/models"
)

type KucoinTrade struct {
	Data struct {
		OrderId string `json:"orderId"`
	} `json:"data"`
}

func (book KucoinTrade) Mapper() any {
	if len(book.Data.OrderId) > 0 {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
