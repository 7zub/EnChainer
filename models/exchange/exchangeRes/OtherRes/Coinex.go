package OtherRes

import (
	"enchainer/models"
)

type CoinexTransfer struct {
	Message string `json:"message"`
}

func (answer CoinexTransfer) Mapper() any {
	if answer.Message == "OK" {
		return models.Result{
			Status: models.OK,
		}
	} else {
		return models.Result{
			Status: models.ERR,
		}
	}
}
