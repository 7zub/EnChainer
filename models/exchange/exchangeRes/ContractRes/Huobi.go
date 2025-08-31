package ContractRes

import (
	"enchainer/models"
)

type HuobiContract struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol       string  `json:"symbol"`
		ContractSize float64 `json:"contract_size"`
	} `json:"data"`
}

func (a HuobiContract) Mapper() any {
	if a.Data[0].ContractSize > 0 {
		return models.Result{
			Status: models.OK,
			Any:    a,
		}
	}

	return models.Result{
		Status: models.ERR,
	}
}
