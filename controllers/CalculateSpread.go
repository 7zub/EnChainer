package controllers

import (
	"awesomeProject/models"
	"fmt"
)

func CalculateSpread(pair models.TradingPair) {
	//pair.OrderBook = append(pair.OrderBook, ApiGetBook("SOLUSDT")) // TODO
	//pair.OrderBook = append(pair.OrderBook, ApiGetBook(pair.Currency))

	ord := ApiGetBook("SOLUSDT")
	omap := BookMapper(ord)
	fmt.Println(omap)
	//jsonBytes, err := json.Marshal(&pair)
	//file, err := os.Create("export.json")
	//
	//if err != nil {
	//	fmt.Println("Unable to create file:", err)
	//	os.Exit(1)
	//}
	//defer file.Close()
	//file.WriteString(string(jsonBytes))

	//a1 := pair.OrderBook[0].Bids

	//if s, err := strconv.ParseFloat(a1[0][0], 64); err == nil {
	//	fmt.Println(s) // 3.14159265
	//}

	//sort.Ints(a1[0][0])
	// Вывод: 15
	//fmt.Println(strconv.ParseFloat(a1[0][0], 64))

	//m, _ := max(strconv.Atoi()

	//for i, order := range pair.OrderBook {
	//	if order.Bids[0][0] == "1" {
	//		//return i, models.Result{Status: "OK"}
	//	}
	//}
}
