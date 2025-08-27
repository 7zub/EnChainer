package ContractRes

import (
	"enchainer/models"
)

type GateioContract []struct {
	Ccy string `json:"name"`
	Cct string `json:"quanto_multiplier"`
}

func (a GateioContract) Mapper() any {
	if len(a) > 0 {
		return models.Result{
			Status: models.OK,
			Any:    a,
		}
	}

	return models.Result{
		Status: models.ERR,
	}
}
