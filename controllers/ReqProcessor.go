package controllers

import (
	"awesomeProject/models"
	"awesomeProject/models/exchangeReq"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ReqCreator(pair models.TradingPair) models.Request {
	for i := 1; i <= 1; i++ {
		switch i {
		case models.BINANCE:
			return models.Request{
				Url:      "exchangeReq.BookRequest",
				Currency: pair.Currency,
				Params:   exchangeReq.BinanceParams{pair.Currency, 10},
			}

		case models.GATEIO:
			return models.Request{
				Url:      "exchangeReq.BookRequest",
				Currency: pair.Currency,
				Params:   exchangeReq.GateioParams{pair.Currency},
			}
		}

	}
	return models.Request{}
}

func UrlCreator(req models.Request) {
	//fields := reflect.TypeOf(req.Params)
	//values := reflect.ValueOf(req.Params)

	rq, err := http.NewRequest("GET", "https://api.binance.com/api/v3/depth", nil)
	if err != nil {
		panic(err)
	}

	//for i := 0; i < fields.NumField(); i++ {
	//	rq.URL.Query().Add(fields.Field(i).Name, "test")
	//}

	q := rq.URL.Query()
	q.Add("symbol", "SOLUSDT")
	//q.Add("param2", "value2")
	rq.URL.RawQuery = q.Encode()

	fmt.Printf("Полный URL: %s\n", rq.URL.String())
	client := http.Client{}
	resp, err := client.Do(rq)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println(result)
}
