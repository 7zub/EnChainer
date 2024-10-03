package controls

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeReq"
)

func BooksPair(pair models.TradingPair) models.Result {
	RqList := []models.IParams{
		exchangeReq.BinanceBookParams{},
		exchangeReq.GateioBookParams{},
	}

	go TaskCreate(&pair, RqList)

	return models.Result{"OK", "Мониторинг пары запущен"}
}
