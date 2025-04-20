package TradeRes

import (
	"enchainer/models"
)

type OkxTrade struct {
	Data []struct {
		OrderId string `json:"ordId"`
	} `json:"data"`
}

func (book OkxTrade) Mapper() any {
	if len(book.Data) > 0 && len(book.Data[0].OrderId) > 0 {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
