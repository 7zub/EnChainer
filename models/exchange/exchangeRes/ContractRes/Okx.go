package ContractRes

import (
	"enchainer/models"
	"strconv"
)

type OkxContract struct {
	Status string `json:"code"`
	Data   []struct {
		Symbol       string `json:"ctValCcy"`
		ContractSize string `json:"ctVal"`
	} `json:"data"`
}

func (a OkxContract) Mapper() any {
	if len(a.Data[0].ContractSize) > 0 {
		cct, _ := strconv.ParseFloat(a.Data[0].ContractSize, 64)

		return models.Result{
			Status: models.OK,
			Any:    cct,
		}
	}

	return models.Result{
		Status: models.ERR,
	}
}
