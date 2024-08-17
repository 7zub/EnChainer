package controllers

import (
	"awesomeProject/models"
	"awesomeProject/models/exchangeResponse"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func ApiGetBook(currency string) exchangeResponse.BinanceBook {
	resp, err := http.Get("https://api.binance.com/api/v3/depth?symbol=" + currency + "&limit=10")

	if err != nil {
		fmt.Println(models.Result{"ERR", "Не удалось подключиться к хосту"})
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var tp exchangeResponse.BinanceBook
		json.Unmarshal(body, &tp)
		return tp
	}
	return exchangeResponse.BinanceBook{}
}

func ApiGetBook1(currency string) exchangeResponse.GateioBook {
	resp, err := http.Get("https://api.gateio.ws/api/v4/spot/order_book?currency_pair=" + currency)

	if err != nil {
		fmt.Println(models.Result{"ERR", "Не удалось подключиться к хосту"})
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var tp exchangeResponse.GateioBook
		json.Unmarshal(body, &tp)
		return tp
	}
	return exchangeResponse.GateioBook{}
}
