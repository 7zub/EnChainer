package ContractReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/ContractRes"
)

type HuobiContractParams struct {
	Ccy string `url:"contract_code"`
}

func (HuobiContractParams) GetParams(any) *models.Request {
	return &models.Request{
		Url:      "https://api.hbdm.com/linear-swap-api/v1/swap_contract_info",
		ReqType:  models.ReqType.Contract,
		Params:   HuobiContractParams{},
		Response: &ContractRes.HuobiContract{},
	}
}
