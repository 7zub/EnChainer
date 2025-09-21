package ContractRes

import (
	"enchainer/models"
	"strconv"
)

type GateioContract []struct {
	Ccy string `json:"name"`
	Cct string `json:"quanto_multiplier"`
}

func (a GateioContract) Mapper() any {
	if len(a) > 0 {
		m := make(map[string]float64)
		for i := range a {
			m[a[i].Ccy], _ = strconv.ParseFloat(a[i].Cct, 64)
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
