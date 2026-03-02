package controls

import (
	"context"
	"enchainer/controls/load"
	"enchainer/models"
	"enchainer/models/exchange/exchangeReq/BookReq"
	"fmt"
	"sync"
	"time"
)

func StartPair(pair *models.TradePair) {
	RqList := []models.IParams{
		BookReq.BinanceBookParams{},
		BookReq.GateioBookParams{},
		BookReq.HuobiBookParams{},
		BookReq.OkxBookParams{},
		BookReq.MexcBookParams{},
		BookReq.BybitBookParams{},
		BookReq.KucoinBookParams{},
		BookReq.CoinexBookParams{},
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
			load.ToLog(models.Result{Status: models.WAR, Message: "Остановлена пара " + pair.Ccy.Currency})
			return
		}
	}
}

func TaskCreate(pair *models.TradePair, reqList []models.IParams) {
	var Wg sync.WaitGroup
	var taskId string

	for _, req := range reqList {
		if SearchReqBlock(pair, GetEx(req)) != nil {
			load.ToLog(models.Result{Status: models.INFO, Message: "Запрос в блок-листе " + pair.Ccy.Currency + " - " + string(GetEx(req))})
			continue
		}

		Wg.Add(1)

		go func(rr models.IParams) {
			ctx, cancel := context.WithTimeout(context.Background(), pair.SessTime-100)
			date, rid := models.GenDescRequest()
			defer Wg.Done()
			defer cancel()
			defer exceptTask(rid)

			rq := rr.GetParams(pair)
			rq.DescRequest(date, rid)
			rq.SendRequest()
			load.ToLog(rq.Log)
			rs := rq.Response.Mapper().(models.OrderBook)

			if isDone(ctx) {
				rq.Log = models.Result{Status: models.WAR, Message: "Задержка запроса " + rq.ReqId + ": " + rq.Url}
				ChanAny <- rq
				load.ToLog(rq.Log)
				return
			}

			if rs.BookExist() {
				rs.ReqId = rq.ReqId
				rs.CreateDate = time.Now()
				rs.TpId = pair.Id
				pair.Mu.Lock()
				pair.OrderBook = append(pair.OrderBook, rs)
				pair.Mu.Unlock()
			} else {
				rb := CreateReqBlock(*rq, pair, rs.Exchange)
				ChanAny <- rb
				ChanAny <- rq

				if rq.Log.Status == models.INFO {
					rq.Log = models.Result{Status: models.WAR, Message: "Некорректный результат запроса " + rq.ReqId}
					load.ToLog(rq.Log)
				}
			}
		}(req)
	}

	Wg.Wait()

	pair.Mu.Lock()
	if len(pair.OrderBook) == 0 {
		pair.Mu.Unlock()
		return
	}
	ob := append([]models.OrderBook(nil), pair.OrderBook...)
	pair.OrderBook = nil
	pair.Mu.Unlock()

	if len(ob) > 1 {
		models.SortOrderBooks(&ob)
		ask, deepAsk := models.GetVolume(&ob[len(ob)-1].Asks)
		bid, deepBid := models.GetVolume(&ob[0].Bids)

		taskId = GenTaskId()
		task := models.TradeTask{
			TaskId: taskId,
			Ccy: models.Ccy{
				Currency:  pair.Ccy.Currency,
				Currency2: pair.Ccy.Currency2,
			},
			Buy: models.Operation{
				Ex:     ob[len(ob)-1].Exchange,
				Price:  ask.Price,
				Volume: ask.Volume,
				Side:   models.Buy,
				Deep:   deepAsk,
				Market: pair.Market,
			},
			Sell: models.Operation{
				Ex:     ob[0].Exchange,
				Price:  bid.Price,
				Volume: bid.Volume,
				Side:   models.Sell,
				Deep:   deepBid,
				Market: pair.Market,
			},
			Spread:     Round((bid.Price/ask.Price-1)*100, 4),
			CreateDate: time.Now(),
			Stage:      models.Creation,
			Status:     models.Done,
		}

		go func(obCopy []models.OrderBook, t models.TradeTask) {
			TradeTask.Store(t.TaskId, &t)
			ChanAny <- t
			PendingHandler(t.Ccy, obCopy)
			TradeTaskHandler(&t)
		}(ob, task)
	}

	for i := range ob {
		ob[i].TaskId = taskId
	}

	go func(obCopy []models.OrderBook) {
		ChanBook <- obCopy
	}(ob)
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func TaskPause() {
	//time.Sleep(15 * time.Second)
	for i, _ := range TradePair {
		if TradePair[i].Status == models.StatusPair.On && SearchPendTask(TradePair[i].Ccy) == nil {
			if TradePair[i].StopCh != nil {
				close(TradePair[i].StopCh)
			}
		}
	}
	load.ToLog(models.Result{Status: models.WAR, Message: fmt.Sprintf("Отключены все пары, кроме pending")})
}

func TaskTime(ccy models.Ccy, sec time.Duration) {
	for i, pair := range TradePair {
		if pair.Ccy == ccy && pair.Status == models.StatusPair.On {
			TradePair[i].SessTime = sec * time.Second
			if TradePair[i].StopCh != nil {
				close(TradePair[i].StopCh)
			}
			StartPair(&TradePair[i])
			load.ToLog(models.Result{
				Status:  models.WAR,
				Message: fmt.Sprintf("Выставлен интервал для %s: %d", ccy.Currency, TradePair[i].SessTime),
			})
		}
	}
}
