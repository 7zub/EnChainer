package OtherRes

import (
	"enchainer/models"
)

type CoinexTransfer struct {
	Message string `json:"message"`
}

func (a CoinexTransfer) Mapper() any {
	if a.Message == "OK" {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
