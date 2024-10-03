package controls

import (
	"awesomeProject/models"
	"fmt"
	"math/rand"
	"time"
)

func TaskCreate(pair *models.TradingPair, reqList []models.IParams) {
	ticker := time.NewTicker(pair.SessTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			start := time.Now()

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
				fmt.Printf("%#v\n", task)

				pair.OrderBook = []models.OrderBook{}
				//ticker.Stop()
			}

			duration := time.Since(start)
			fmt.Printf("Время выполнения: %v\n. Ожидание следующего интервала...\n\n", duration)

			for _, req := range reqList {
				go func(rr models.IParams) {
					r := rr.GetParams(pair.Ccy)
					r.SendRequest()

					pair.OrderBook = append(pair.OrderBook, r.Response.Mapper())
				}(req)
			}

			//case <-stop:
			//	// Получен сигнал об остановке
			//	fmt.Println("Остановлен")
			//	return
		}
	}
}
