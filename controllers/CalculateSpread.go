package controllers

import (
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"math"
	"os"
)

func Calc__TaskManDev(pair models.TradingPair) models.Result { //TODO
	//UrlCreator(ReqCreator(pair))
	start1()

	return models.Result{"test", "em"}
}

func CalculateSpread(pair models.TradingPair) models.Result {
	pair.OrderBook = append(pair.OrderBook, ApiGetBook("SOLUSDT").BookMapper()) // TODO
	pair.OrderBook = append(pair.OrderBook, ApiGetBook(pair.Currency).BookMapper())

	jsonBytes, err := json.Marshal(&pair)
	file, err := os.Create("export.json")

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
					Currency:     pair.Currency,
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
