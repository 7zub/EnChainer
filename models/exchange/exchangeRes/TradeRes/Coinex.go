package TradeRes

import (
	"enchainer/models"
)

type CoinexTrade struct {
	Data struct {
		OrderId int64 `json:"order_id"`
	} `json:"data"`
}

func (book CoinexTrade) Mapper() any {
	if book.Data.OrderId > 1 {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
