package controls

import (
	"enchainer/models"
	"sync"
	"time"
)

var ReqBlock sync.Map

func CreateReqBlock(rq models.Request, ccy models.Ccy, ex models.Exchange) *models.RequestBlock {
	var reqb *models.RequestBlock
	v, ok := ReqBlock.Load(ccy.Currency + string(ex))

	if ok {
		reqb = v.(*models.RequestBlock)
		reqb.Active = true
	} else {
		reqb = &models.RequestBlock{
			ReqId:      rq.ReqId,
			Ccy:        ccy,
			Ex:         ex,
			ReasonCode: rq.Code,
			Reason:     rq.ResponseRaw,
			CreateDate: time.Now(),
			Active:     true,
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
			if b.ReasonCode == 400 || time.Since(b.CreateDate) < 110*time.Second {
				res = b.ReqId
			} else {
				b.Active = false
				ReqBlock.Store(ccy.Currency+string(ex), b)
				ChanAny <- b
			}
			return false
		}
		return true
	})
	return res
}
