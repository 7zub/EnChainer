package controls

import (
	"enchainer/models"
	"sync"
)

var ReqBlock sync.Map

func SearchReqBlock(ccy models.Ccy, ex models.Exchange) string {
	var res string
	ReqBlock.Range(func(key, val any) bool {
		b, _ := val.(models.RequestBlock)
		if ccy == b.Ccy && ex == b.Ex {
			res = b.ReqId
			return false
		}
		return true
	})
	return res
}
