package controls

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeReq"
	"fmt"
	"math/rand"
	"time"
)

var pairMap = make(map[string]chan bool)

func BooksPair(pair *models.TradingPair) models.Result {
	RqList := []models.IParams{
		exchangeReq.BinanceBookParams{},
		exchangeReq.GateioBookParams{},
	}

	go TaskTicker(pair, RqList)

	return models.Result{"OK", "Мониторинг пары запущен"}
}

func TaskTicker(pair *models.TradingPair, reqList []models.IParams) {
	stop := make(chan bool)
	pairMap[pair.PairId] = stop
	ticker := time.NewTicker(pair.SessTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			TaskCreate(pair, reqList)

		case <-stop:
			fmt.Println("Остановлен")
			return
		}
	}
}

func TaskCreate(pair *models.TradingPair, reqList []models.IParams) {
	if len(pair.OrderBook) > 0 {
		models.SortOrderBooks(&pair.OrderBook)

		task := models.TradeTask{
			TaskId: rand.Intn(1000),
			Currency: models.Ccy{
				Currency:  pair.Ccy.Currency,
				Currency2: pair.Ccy.Currency2,
			},
			Buy: models.Operation{
				Exchange: pair.OrderBook[len(pair.OrderBook)-1].Exchange,
				Price:    pair.OrderBook[len(pair.OrderBook)-1].Asks[0].Price,
				Volume:   nil,
			},
			Sell: models.Operation{
				Exchange: pair.OrderBook[0].Exchange,
				Price:    pair.OrderBook[0].Bids[0].Price,
				Volume:   nil,
			},
			Profit: pair.OrderBook[0].Bids[0].Price/pair.OrderBook[0].Asks[0].Price - 1,
		}

		TradeTask = append(TradeTask, task)
		go func() {
			SaveBookDb(*&pair)
			SaveTradeDb(&task)
			pair.OrderBook = []models.OrderBook{}
		}()
	}

	for _, req := range reqList {
		go func(rr models.IParams) {
			r := rr.GetParams(pair.Ccy)
			r.SendRequest()
			// todo "сделать принудительное завершение горудин по истечению таймера + try"
			newbook := r.Response.Mapper()
			newbook.ReqDate = r.ReqDate
			pair.OrderBook = append(pair.OrderBook, newbook)
		}(req)
	}
}

func TaskStop(pair string) {
	if quit, exists := pairMap[pair]; exists {
		close(quit)
		delete(pairMap, pair)
		fmt.Printf("Горутина для торговой пары %s остановлена\n", pair)
	} else {
		fmt.Printf("Горутина для торговой пары %s не найдена\n", pair)
	}
}
