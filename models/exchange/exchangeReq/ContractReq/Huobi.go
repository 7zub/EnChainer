package ContractReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/ContractRes"
)

type HuobiContractParams struct {
	Ccy string `url:"contract_code"`
}

func (HuobiContractParams) GetParams(task any) *models.Request {
	t := task.(models.OperationTask)

	return &models.Request{
		Url:      "https://api.hbdm.com/linear-swap-api/v1/swap_contract_info",
		ReqType:  models.ReqType.Contract,
		Params:   HuobiContractParams{Ccy: t.Ccy.Currency + "-" + t.Ccy.Currency2},
		Response: &ContractRes.HuobiContract{},
	}
}
