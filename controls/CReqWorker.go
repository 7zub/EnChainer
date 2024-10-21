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
			if pair.StopCh != nil {
				close(pair.StopCh)
			}
			pair.StopCh = make(chan struct{})
			TaskCreate(pair, reqList)

		case <-stop:
			fmt.Println("Остановлен")
			if pair.StopCh != nil {
				close(pair.StopCh)
			}
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
			SaveBookDb(pair)
			SaveTradeDb(&task)
			pair.OrderBook = []models.OrderBook{}
		}()
	}

	for _, req := range reqList {
		go func(rr models.IParams) {
			done := make(chan struct{})
			go func() {
				r := rr.GetParams(pair.Ccy)
				r.SendRequest()
				go SaveReqDb(r)
				newbook := r.Response.Mapper()
				newbook.ReqDate = r.ReqDate
				//newbook.ReqId = r.ReqId
				pair.OrderBook = append(pair.OrderBook, newbook)
				close(done)
			}()

			select {
			case <-pair.StopCh:
				fmt.Println("Горутина остановлена по сигналу")
				return
			case <-done:
				fmt.Println("Операция завершена")
			case <-time.After(10 * time.Second):
				// Если операция зависла и не завершилась за 10 секунд
				fmt.Println("Операция прервана по тайм-ауту")
				return
			}
		}(req)

		//for {
		//	select {
		//	case <-pair.StopCh:
		//		fmt.Println("Горутина остановлена по сигналу")
		//		return
		//	}
		//}
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
