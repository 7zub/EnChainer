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
	t := task.(models.OperationTask)

	return &models.Request{
		Url:      "https://www.okx.com/api/v5/public/instruments",
		ReqType:  models.ReqType.Contract,
		Params:   OkxContractParams{Ccy: t.Ccy.Currency + "-" + t.Ccy.Currency2, Type: "SWAP"},
		Response: &ContractRes.OkxContract{},
	}
}
