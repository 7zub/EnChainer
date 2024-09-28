package controllers

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeReq"
	"encoding/json"
	"fmt"
	"math"
	"os"
)

func BooksPair(pair models.TradingPair) models.Result {
	RqList := []models.IParams{
		exchangeReq.BinanceBookParams{},
		exchangeReq.GateioBookParams{},
	}

	for _, req := range RqList {
		go TaskCreate(&pair, req)
		//go req.GetParams(pair.Ccy).SendRequest()
	}

	return models.Result{"OK", "Мониторинг пары запущен"}
}

func CalcSpread(pair models.TradingPair) models.Result {
	pair.OrderBook = append(pair.OrderBook, ApiGetBook("SOLUSDT").Mapper())
	//pair.OrderBook = append(pair.OrderBook, ApiGetBook(pair.Currency).Mapper())

	jsonBytes, err := json.Marshal(&pair)
	file, err := os.Create("files/export.json")

	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(string(jsonBytes))

	var result string

	for i, order := range pair.OrderBook {
		for j, orderOther := range pair.OrderBook {
			if i != j && order.Asks[0].Price < orderOther.Bids[0].Price {

				task := models.TradeTask{
					TaskId:       1,
					Currency:     pair.Ccy.Currency,
					ExchangeBuy:  order.Exchange,
					ExchangeSell: orderOther.Exchange,
					PriceBuy:     order.Asks[0].Price,
					PriceSell:    orderOther.Bids[0].Price,
					Profit:       math.Round(orderOther.Bids[0].Price/order.Asks[0].Price - 1),
				}

				str, _ := json.Marshal(task)
				result = result + string(str)
			}
		}
	}

	return models.Result{Status: "OK", Message: result}
}
