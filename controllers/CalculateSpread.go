package controllers

import (
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"os"
)

func CalculateSpread(pair models.TradingPair) {
	pair.OrderBook = append(pair.OrderBook, ApiGetBook("SOLUSDT")) // TODO
	pair.OrderBook = append(pair.OrderBook, ApiGetBook(pair.Currency))

	jsonBytes, err := json.Marshal(&pair)
	file, err := os.Create("export.json")

	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(string(jsonBytes))

	a1 := pair.OrderBook[0].Bids
	//a2, _ := strconv.Atoi(pair.OrderBook[0].Asks[0][0])
	//a2 := a1.(map[string]interface{})
	//ab := aa[0]
	fmt.Println(a1)

	//m, _ := max(strconv.Atoi()

	for i, order := range pair.OrderBook {
		if order.Bids[0][0] == "1" {
			return i, models.Result{Status: "OK"}
		}
	}
}
