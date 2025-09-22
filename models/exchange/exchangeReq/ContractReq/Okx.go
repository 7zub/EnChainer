package ContractReq

import (
	"enchainer/models"
	"enchainer/models/exchange/exchangeRes/ContractRes"
)

type OkxContractParams struct {
	Ccy  string `url:"uly"`
	Type string `url:"instType"`
}

func (OkxContractParams) GetParams(task any) *models.Request {
	return &models.Request{
		Url:      "https://www.okx.com/api/v5/public/instruments",
		ReqType:  models.ReqType.Contract,
		Params:   OkxContractParams{Type: "SWAP"},
		Response: &ContractRes.OkxContract{},
	}
}
