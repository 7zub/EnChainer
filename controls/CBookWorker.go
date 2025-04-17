package controls

import (
	"context"
	"enchainer/models"
	"enchainer/models/exchange/exchangeReq/BookReq"
	"time"
)

func StartPair(pair *models.TradePair) {
	RqList := []models.IParams{
		BookReq.BinanceBookParams{},
		BookReq.GateioBookParams{},
		//BookReq.HuobiBookParams{},
		//BookReq.OkxBookParams{},
		BookReq.MexcBookParams{},
		BookReq.BybitBookParams{},
		//BookReq.KucoinBookParams{},
	}
	go TaskTicker(pair, RqList)
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
			ToLog("Остановлена пара " + pair.Ccy.Currency)
			return
		}
	}
}

func TaskCreate(pair *models.TradePair, reqList []models.IParams) {
	if len(pair.OrderBook) > 1 {
		models.SortOrderBooks(&pair.OrderBook)

		taskId := GenTaskId()
		task := models.TradeTask{
			TaskId: taskId,
			Ccy: models.Ccy{
				Currency:  pair.Ccy.Currency,
				Currency2: pair.Ccy.Currency2,
			},
			Buy: models.Operation{
				Ex:     pair.OrderBook[len(pair.OrderBook)-1].Exchange,
				Price:  pair.OrderBook[len(pair.OrderBook)-1].Asks[0].Price,
				Volume: pair.OrderBook[len(pair.OrderBook)-1].Asks[0].Volume,
				Side:   models.Buy,
			},
			Sell: models.Operation{
				Ex:     pair.OrderBook[0].Exchange,
				Price:  pair.OrderBook[0].Bids[0].Price,
				Volume: pair.OrderBook[0].Bids[0].Volume,
				Side:   models.Sell,
			},
			Spread:     (pair.OrderBook[0].Bids[0].Price/pair.OrderBook[len(pair.OrderBook)-1].Asks[0].Price - 1) * 100,
			CreateDate: time.Now(),
			Stage:      models.Creation,
			Status:     models.Done,
		}
		go func() {
			TradeTask.Store(taskId, task)
			SaveDb(&task)
			TradeTaskHandler(LoadTask(taskId))
		}()

	}

	if len(pair.OrderBook) > 0 {
		go func() {
			SaveDb(pair)
			pair.OrderBook = []models.OrderBook{}
		}()
	}

	for _, req := range reqList {
		if SearchReqBlock(pair.Ccy, GetEx(req)) != "" {
			ToLog(models.Result{Status: models.INFO, Message: "Запрос в блок-листе " + pair.Ccy.Currency + " - " + string(GetEx(req))})
			continue
		}

		go func(rr models.IParams) {
			ctx, cancel := context.WithTimeout(context.Background(), pair.SessTime-100)
			date, rid := models.GenDescRequest()
			defer cancel()
			defer exceptTask(rid)

			rq := rr.GetParams(pair.Ccy)
			rq.DescRequest(date, rid)
			rq.SendRequest()
			ToLog(*rq)
			go SaveDb(rq)
			rs := rq.Response.Mapper().(models.OrderBook)

			if isDone(ctx) {
				rq.Log = models.Result{Status: models.WAR, Message: "Задержка запроса " + rq.ReqId + ": " + rq.Url}
				ToLog(*rq)
				return
			}

			if rs.BookExist() {
				rs.ReqId = rq.ReqId
				pair.OrderBook = append(pair.OrderBook, rs)
				go PendingHandler(pair.Ccy, rs)
			} else {
				rb := CreateReqBlock(rq.ReqId, pair.Ccy, rs.Exchange)
				SaveDb(&rb)

				if rq.Log.Status == models.INFO {
					rq.Log = models.Result{Status: models.WAR, Message: "Некорректный результат запроса " + rq.ReqId}
					ToLog(*rq)
				}
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
