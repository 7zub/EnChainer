package controls

import (
	"context"
	"enchainer/models"
	"enchainer/models/exchange/exchangeReq/BookReq"
	"fmt"
	"time"
)

func BooksPair(pair *models.TradePair) {
	RqList := []models.IParams[models.Ccy]{
		BookReq.BinanceBookParams{},
		BookReq.GateioBookParams{},
		BookReq.HuobiBookParams{},
		BookReq.OkxBookParams{},
		BookReq.MexcBookParams{},
		BookReq.BybitBookParams{},
		BookReq.KucoinBookParams{},
	}
	go TaskTicker(pair, RqList)
}

func TaskTicker(pair *models.TradePair, reqList []models.IParams[models.Ccy]) {
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

func TaskCreate(pair *models.TradePair, reqList []models.IParams[models.Ccy]) {
	if len(pair.OrderBook) > 0 {
		models.SortOrderBooks(&pair.OrderBook)

		task := models.TradeTask{
			Ccy: models.Ccy{
				Currency:  pair.Ccy.Currency,
				Currency2: pair.Ccy.Currency2,
			},
			Buy: models.Operation{
				Ex:     pair.OrderBook[len(pair.OrderBook)-1].Exchange,
				Price:  pair.OrderBook[len(pair.OrderBook)-1].Asks[0].Price,
				Volume: pair.OrderBook[len(pair.OrderBook)-1].Asks[0].Volume,
			},
			Sell: models.Operation{
				Ex:     pair.OrderBook[0].Exchange,
				Price:  pair.OrderBook[0].Bids[0].Price,
				Volume: pair.OrderBook[0].Bids[0].Volume,
			},
			Spread: (pair.OrderBook[0].Bids[0].Price/pair.OrderBook[len(pair.OrderBook)-1].Asks[0].Price - 1) * 100,
		}

		TradeTask = append(TradeTask, task)
		go func() {
			SaveBookDb(pair)
			SaveTradeDb(&task)
			pair.OrderBook = []models.OrderBook{}
		}()
	}

	for _, req := range reqList {
		go func(rr models.IParams[models.Ccy]) {
			ctx, cancel := context.WithTimeout(context.Background(), pair.SessTime-100)
			date, rid := models.GenDescRequest()
			defer cancel()
			defer exceptTask(rid)

			rq := rr.GetParams(pair.Ccy)
			rq.DescRequest(date, rid)
			rq.SendRequest()
			ToLog(*rq)
			go SaveReqDb(rq)
			rs := rq.Response.Mapper().(models.OrderBook)

			if isDone(ctx) {
				rq.Log = models.Result{Status: models.WAR, Message: "Задержка запроса " + rq.ReqId + ": " + rq.Url}
				ToLog(*rq)
				return
			}

			if rs.BookExist() {
				rs.ReqDate = rq.ReqDate
				rs.ReqId = rq.ReqId
				pair.OrderBook = append(pair.OrderBook, rs)
			} else {
				rq.Log = models.Result{Status: models.WAR, Message: "Некорректный результат запроса " + rq.ReqId}
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
