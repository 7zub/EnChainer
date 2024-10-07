package controls

import (
	"awesomeProject/models"
	"encoding/json"
	"net/http"
	"slices"
	"time"
)

var TradingPair = []models.TradingPair{}

func BookControl(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(TradingPair)
}

func AddPair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	tp := models.TradingPair{
		PairId: "P_" + params.Get("currency"),
		Name:   params.Get("name"),
		Desc:   params.Get("desc"),
		Ccy: models.Ccy{
			Currency:  params.Get("currency"),
			Currency2: "USDT",
		},
		Status:     models.Off,
		SessTime:   2 * time.Second,
		CreateDate: time.Now(),
	}

	TradingPair = append(TradingPair, tp)
	db1(&tp)
	//db1(tp)

	json.NewEncoder(w).Encode(models.Result{"OK", "Добавлена пара: " + params.Get("currency")})
}

func DeletePair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	var delid string

	for i, pair := range TradingPair {
		if pair.PairId == params.Get("id") {
			TradingPair = slices.Delete(TradingPair, i, i+1)
			delid = pair.PairId
			break
		}
	}

	if delid != "" {
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

func SearchPair(pairid string) (int, models.Result) {
	for i, pair := range TradingPair {
		if pairid == pair.PairId {
			return i, models.Result{Status: "OK"}
		}
	}
	return -1, models.Result{"ERR", "Пары не существует"}
}
