package controls

import (
	"enchainer/models"
	"sync"
	"time"
)

var ReqBlock sync.Map

func CreateReqBlock(rq models.Request, ccy models.Ccy, ex models.Exchange) *models.RequestBlock {
	var reqb *models.RequestBlock
	v, ok := ReqBlock.Load(ccy.Currency + string(ex)) //TODO

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
			RepeatDate: time.Now(),
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
			loc, _ := time.LoadLocation("Europe/Moscow")
			rptDate := time.Date(b.RepeatDate.Year(), b.RepeatDate.Month(), b.RepeatDate.Day(), b.RepeatDate.Hour(), b.RepeatDate.Minute(), b.RepeatDate.Second(), b.RepeatDate.Nanosecond(), loc)

			if b.ReasonCode == 400 || time.Since(rptDate) < models.Const.TimeoutBlock*time.Second {
				res = b.ReqId
			} else {
				b.Active = false
				b.RepeatDate = time.Now()
				ReqBlock.Store(ccy.Currency+string(ex), b)
				ChanAny <- b
			}
			return false
		}
		return true
	})
	return res
}
