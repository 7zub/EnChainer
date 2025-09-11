package controls

import (
	"enchainer/models"
	"sync"
	"time"
)

var ReqBlock sync.Map

func CreateReqBlock(rq models.Request, pair *models.TradePair, ex models.Exchange) *models.RequestBlock {
	var reqb *models.RequestBlock
	v, ok := ReqBlock.Load(string(pair.Market) + pair.Ccy.Currency + string(ex))

	if ok {
		reqb = v.(*models.RequestBlock)
		reqb.Active = true
	} else {
		reqb = &models.RequestBlock{
			ReqId:      rq.ReqId,
			Market:     pair.Market,
			Ccy:        pair.Ccy,
			Ex:         ex,
			ReasonCode: rq.Code,
			Reason:     rq.ResponseRaw,
			CreateDate: time.Now(),
			RepeatDate: time.Now(),
			Active:     true,
		}
	}

	ReqBlock.Store(string(pair.Market)+pair.Ccy.Currency+string(ex), reqb)
	return reqb
}

func SearchReqBlock(pair *models.TradePair, ex models.Exchange) *string {
	v, _ := ReqBlock.Load(string(pair.Market) + pair.Ccy.Currency + string(ex))
	b, _ := v.(*models.RequestBlock)

	if b != nil && b.Active == true {
		loc, _ := time.LoadLocation("Europe/Moscow")
		rptDate := time.Date(b.RepeatDate.Year(), b.RepeatDate.Month(), b.RepeatDate.Day(), b.RepeatDate.Hour(), b.RepeatDate.Minute(), b.RepeatDate.Second(), b.RepeatDate.Nanosecond(), loc)

		if b.ReasonCode == 400 || time.Since(rptDate) < models.Const.TimeoutBlock*time.Second {
			return &b.ReqId
		} else {
			b.Active = false
			b.RepeatDate = time.Now()
			ReqBlock.Store(string(pair.Market)+pair.Ccy.Currency+string(ex), b)
			ChanAny <- b
		}
	}
	return nil
}
