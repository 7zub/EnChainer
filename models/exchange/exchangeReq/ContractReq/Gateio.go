package ContractReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/ContractRes"
)

type GateioContractParams struct {
	Ccy      string
	Contract float64
}

func (GateioContractParams) GetParams(any) *models.Request {
	return &models.Request{
		Url:      "https://api.gateio.ws/api/v4/futures/usdt/contracts",
		ReqType:  models.ReqType.Contract,
		Params:   GateioContractParams{},
		Response: &ContractRes.GateioContract{},
	}
}
