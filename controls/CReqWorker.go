package controls

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeReq"
	"context"
	"fmt"
	"math/rand"
	"time"
)

func BooksPair(pair *models.TradingPair) models.Result {
	RqList := []models.IParams{
		exchangeReq.BinanceBookParams{},
		exchangeReq.GateioBookParams{},
		exchangeReq.HuobiBookParams{},
		//exchangeReq.OkxBookParams{},
	}

	go TaskTicker(pair, RqList)

	return models.Result{"OK", "Мониторинг пары запущен"}
}

func TaskTicker(pair *models.TradingPair, reqList []models.IParams) {
	pair.StopCh = make(chan struct{})
	ticker := time.NewTicker(pair.SessTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			TaskCreate(pair, reqList)

		case <-pair.StopCh:
			fmt.Println("Остановлена пара TaskTicker")
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
			Profit: (pair.OrderBook[0].Bids[0].Price/pair.OrderBook[len(pair.OrderBook)-1].Asks[0].Price - 1) * 100,
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			rq := rr.GetParams(pair.Ccy)
			rq.DescRequest()
			go SaveReqDb(rq)
			rq.SendRequest()
			rs := rq.Response.Mapper()

			if isDone(ctx) {
				fmt.Println("isDone()")
				return
			}

			if rs.BookExist() {
				rs.ReqDate = rq.ReqDate
				rs.ReqId = rq.ReqId
				pair.OrderBook = append(pair.OrderBook, rs)
			}
		}(req)
	}
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
