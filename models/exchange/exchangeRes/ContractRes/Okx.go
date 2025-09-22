package ContractRes

import (
	"enchainer/models"
	"strconv"
	"strings"
)

type OkxContract struct {
	Status string `json:"code"`
	Data   []struct {
		Ccy  string `json:"uly"`
		Cct  string `json:"ctVal"`
		Type string `json:"ctType"`
	} `json:"data"`
}

func (a OkxContract) Mapper() any {
	if len(a.Data) > 0 {
		m := make(map[string]float64)
		for i := range a.Data {
			if a.Data[i].Type == "linear" {
				m[strings.Replace(a.Data[i].Ccy, "-", "_", 1)], _ = strconv.ParseFloat(a.Data[i].Cct, 64)
			}
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
