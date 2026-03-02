package controls

import (
	"enchainer/controls/load"
	"enchainer/models"
	"fmt"
	"time"
)

var PairInfo = make(map[string]*models.TradePairInfo)

func LoadCcyInfo() {
	LoadPairInfoDb()

	var t time.Time
	for _, p := range PairInfo {
		t = UniZone(p.ReloadDate)
		break
	}

	if time.Since(t) < models.Const.TimeoutCcyInfo {
		load.ToLog(models.Result{
			Status:  models.WAR,
			Message: fmt.Sprintf("Обновление справочника контрактов не требуется, не прошло %.2f часов", models.Const.TimeoutCcyInfo.Hours())})
		return
	} else {
		load.ToLog(models.Result{
			Status:  models.WAR,
			Message: fmt.Sprintf("Перезагрузка справочника контрактов: текущий размер %d", len(PairInfo))})
	}

	var exList = []models.Exchange{
		models.GATEIO,
		models.HUOBI,
		models.OKX,
	}

	for _, ex := range exList {
		var pairInfo, _ = CreateAction(models.Operation{Ex: ex}, models.ReqType.Contract)
		data, _ := pairInfo.Any.(map[string]float64)

		for _, pair := range TradePair {
			if pair.Status == models.StatusPair.On && data[pair.Ccy.Currency+"_"+pair.Ccy.Currency2] != 0 {
				if c, ok := PairInfo[pair.Ccy.Currency+"-"+string(ex)]; ok {
					if c.Cct != data[pair.Ccy.Currency+"_"+pair.Ccy.Currency2] {
						cctOld := c.Cct
						c.Cct = data[pair.Ccy.Currency+"_"+pair.Ccy.Currency2]
						c.UpdateDate = &[]time.Time{time.Now()}[0]
						load.ToLog(models.Result{
							Status: models.WAR,
							Message: fmt.Sprintf("Обновился контракт: %f на %f, Ex: %s, Ccy: %s",
								cctOld, c.Cct, c.Ex, c.Ccy.Currency)})
					}
					c.ReloadDate = time.Now()
					SaveDb(c)
				} else {
					PairInfo[pair.Ccy.Currency+"-"+string(ex)] = &models.TradePairInfo{
						Ex:  ex,
						Ccy: pair.Ccy,
						Cct: data[pair.Ccy.Currency+"_"+pair.Ccy.Currency2],
					}
					SaveDb(PairInfo[pair.Ccy.Currency+"-"+string(ex)])
				}
			}
		}
	}
}
