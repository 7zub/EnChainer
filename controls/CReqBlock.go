package controls

import (
	"enchainer/models"
	"sync"
	"time"
)

var ReqBlock sync.Map

func CreateReqBlock(rid string, ccy models.Ccy, ex models.Exchange) *models.RequestBlock {
	var reqb *models.RequestBlock
	v, ok := ReqBlock.Load(ccy.Currency + string(ex))

	if ok {
		reqb = v.(*models.RequestBlock)
		reqb.Active = true
	} else {
		reqb = &models.RequestBlock{
			ReqId:  rid,
			Ccy:    ccy,
			Ex:     ex,
			Active: true,
		}
	}

	ReqBlock.Store(ccy.Currency+string(ex), reqb)
	return reqb

}

func SearchReqBlock(ccy models.Ccy, ex models.Exchange) string {
	var res string
	ReqBlock.Range(func(key, val any) bool {
		b, _ := val.(*models.RequestBlock)
		if ccy == b.Ccy && ex == b.Ex && b.Active == true {
			if time.Since(b.CreateDate) > 50*time.Second {
				b.Active = false
				ReqBlock.Store(ccy.Currency+string(ex), b)
			} else {
				res = b.ReqId
			}
			return false
		}
		return true
	})
	return res
}
