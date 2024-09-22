package controllers

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange/exchangeReq"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

func start1() {
	r := []models.IParams{
		exchangeReq.BinanceBookParams{},
		//exchangeReq.GateioBookParams{},
	}

	for _, i := range r {
		f := i.GetParams("NEO")
		println(f.Params)
	}
}

func UrlCreator(request models.Request) http.Request {
	fields := reflect.TypeOf(request.Params)
	values := reflect.ValueOf(request.Params)

	rq, err := http.NewRequest("GET", request.Url, nil)
	if err != nil {
		panic(err)
	}

	for i := 0; i < fields.NumField(); i++ {
		rq.URL.Query().Add(fields.Field(i).Name, values.Field(1).String())
	}

	rq.URL.RawQuery = rq.URL.Query().Encode()

	fmt.Printf("Полный URL: %s\n", rq.URL.String())
	client := http.Client{}
	resp, err := client.Do(rq)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println(result)

	return *rq
}
