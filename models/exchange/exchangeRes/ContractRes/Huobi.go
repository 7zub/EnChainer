package ContractRes

import (
	"enchainer/models"
	"strings"
)

type HuobiContract struct {
	Status string `json:"status"`
	Data   []struct {
		Ccy string  `json:"contract_code"`
		Cct float64 `json:"contract_size"`
	} `json:"data"`
}

func (a HuobiContract) Mapper() any {
	if len(a.Data) > 0 {
		m := make(map[string]float64)
		for i := range a.Data {
			m[strings.Replace(a.Data[i].Ccy, "-", "_", 1)] = a.Data[i].Cct
		}

		return models.Result{
			Status: models.OK,
			Any:    m,
		}
	}

	return models.Result{
		Status: models.ERR,
	}
}
