package controllers

import (
	"awesomeProject/models"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"time"
)

var TradingPair = []models.TradingPair{}

func Monitor(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(TradingPair)
}

func Book(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	i, res := SearchPair(params.Get("id"))

	if i != -1 {
		//go req.GetParams(pair.Ccy).SendRequest()
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func AddPair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	Id := 1 //rand.Intn(10000)

	TradingPair = append(TradingPair, models.TradingPair{
		Id:   Id,
		Name: params.Get("name"),
		Desc: params.Get("desc"),
		Ccy: models.Ccy{
			Currency:  params.Get("currency"),
			Currency2: "USDT",
		},
		Status:   models.Off,
		SessTime: 5 * time.Second,
	})

	json.NewEncoder(w).Encode(models.Result{"OK", "Добавлена пара: " + params.Get("currency")})
}

func DeletePair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var delid int

	for i, pair := range TradingPair {
		if strconv.Itoa(pair.Id) == params.Get("id") {
			TradingPair = slices.Delete(TradingPair, i, i+1)
			delid = pair.Id
			break
		}
	}

	if delid > 0 {
		json.NewEncoder(w).Encode(models.Result{"OK", "Пара удалена"})
	} else {
		json.NewEncoder(w).Encode(models.Result{"ERR", "Пары не существует"})
	}
}

func OnPair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	i, res := SearchPair(params.Get("id"))

	if i != -1 {
		TradingPair[i].Status = models.On
		json.NewEncoder(w).Encode(BooksPair(TradingPair[i]))
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func OffPair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	i, res := SearchPair(params.Get("id"))

	if i != -1 {
		TradingPair[i].Status = models.Off
	}

	json.NewEncoder(w).Encode(res)
}

func SearchPair(id string) (int, models.Result) {
	cid, _ := strconv.Atoi(id)

	for i, pair := range TradingPair {
		if cid == pair.Id {
			return i, models.Result{Status: "OK"}
		}
	}
	return -1, models.Result{"ERR", "Пары не существует"}
}
