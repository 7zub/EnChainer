package controls

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeReq"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func BooksPair(pair *models.TradePair) models.Result {
	RqList := []models.IParams{
		exchangeReq.BinanceBookParams{},
		exchangeReq.GateioBookParams{},
		exchangeReq.HuobiBookParams{},
		exchangeReq.OkxBookParams{},
		exchangeReq.MexcBookParams{},
		exchangeReq.BybitBookParams{},
		exchangeReq.KucoinBookParams{},
	}

	go TaskTicker(pair, RqList)

	return models.Result{"OK", "Мониторинг пары запущен"}
}

func TaskTicker(pair *models.TradePair, reqList []models.IParams) {
	pair.StopCh = make(chan struct{})
	ticker := time.NewTicker(pair.SessTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			TaskCreate(pair, reqList)

		case <-pair.StopCh:
			fmt.Println("Остановлена пара " + pair.Ccy.Currency)
			return
		}
	}
}

func TaskCreate(pair *models.TradePair, reqList []models.IParams) {
	if len(pair.OrderBook) > 0 {
		models.SortOrderBooks(&pair.OrderBook)

		task := models.TradeTask{
			TaskId: rand.Intn(1000),
			Ccy: models.Ccy{
				Currency:  pair.Ccy.Currency,
				Currency2: pair.Ccy.Currency2,
			},
			Buy: models.Operation{
				Ex:     pair.OrderBook[len(pair.OrderBook)-1].Exchange,
				Price:  pair.OrderBook[len(pair.OrderBook)-1].Asks[0].Price,
				Volume: nil,
			},
			Sell: models.Operation{
				Ex:     pair.OrderBook[0].Exchange,
				Price:  pair.OrderBook[0].Bids[0].Price,
				Volume: nil,
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
			defer exceptTask()

			rq := rr.GetParams(pair.Ccy)
			rq.DescRequest()
			go SaveReqDb(rq)
			rq.SendRequest()
			ToLog(*rq)
			go SaveReqDb(rq)
			rs := rq.Response.Mapper()

			if isDone(ctx) {
				log.Println("isDone()")
				return
			}

			if rs.BookExist() {
				rs.ReqDate = rq.ReqDate
				rs.ReqId = rq.ReqId
				pair.OrderBook = append(pair.OrderBook, rs)
			} else {
				rq.Log = models.Result{Status: models.ERR, Message: "Некорректный результат запроса"}
				ToLog(*rq)
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
