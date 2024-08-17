package controllers

import (
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func CalculateSpread(pair models.TradingPair) models.Result {
	pair.OrderBook = append(pair.OrderBook, ApiGetBook("SOLUSDT").BookMapper()) // TODO
	pair.OrderBook = append(pair.OrderBook, ApiGetBook(pair.Currency).BookMapper())
	pair.OrderBook = append(pair.OrderBook, ApiGetBook1("NEO_USDT").BookMapper())

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
				s := fmt.Sprintf("%.6f", order.Asks[0].Price)
				o := fmt.Sprintf("%.6f", orderOther.Bids[0].Price)

				result = result +
					"[Купить в " + strconv.Itoa(order.Exchange) + " по " + s + ", продать в " + strconv.Itoa(orderOther.Exchange) + " по " + o + "]\n"
			}
		}
	}

	return models.Result{Status: "OK", Message: result}
}
