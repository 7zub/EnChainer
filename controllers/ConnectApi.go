package controllers

import (
	"awesomeProject/models"
	exchangeResp2 "awesomeProject/models/exchange/exchangeRes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func ApiGetBook(currency string) exchangeResp2.BinanceBook {
	resp, err := http.Get("https://api.binance.com/api/v3/depth?symbol=" + currency + "&limit=10")

	if err != nil {
		fmt.Println(models.Result{"ERR", "Не удалось подключиться к хосту"})
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var tp exchangeResp2.BinanceBook
		json.Unmarshal(body, &tp)
		return tp
	}
	return exchangeResp2.BinanceBook{}
}
