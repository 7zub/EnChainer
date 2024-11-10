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
		Title:  params.Get("title"),
		Ccy: models.Ccy{
			Currency:  params.Get("currency"),
			Currency2: "USDT",
		},
		Status:   models.Off,
		SessTime: 2000 * time.Millisecond,
	}

	if i, _ := SearchPair(params.Get("id")); i == -1 {

		TradingPair = append(TradingPair, tp)
		SaveBookDb(&TradingPair[len(TradingPair)-1])

		json.NewEncoder(w).Encode(models.Result{"OK", "Пара " + params.Get("currency") + " добавлена"})
	} else {
		json.NewEncoder(w).Encode(models.Result{"ERR", "Пара " + params.Get("currency") + " уже существует"})
	}
}

func DeletePair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if i, res := SearchPair(params.Get("id")); i != -1 {
		if TradingPair[i].Status == models.Off {
			DeleteBookDb(&TradingPair[i])
			TradingPair = slices.Delete(TradingPair, i, i+1)

			json.NewEncoder(w).Encode(models.Result{"OK", "Пара удалена"})
		} else {
			json.NewEncoder(w).Encode(models.Result{"ERR", "Пара должна быть остановлена"})
		}
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func OnPair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if i, res := SearchPair(params.Get("id")); i != -1 {
		if TradingPair[i].Status != models.On {
			TradingPair[i].Status = models.On
			json.NewEncoder(w).Encode(BooksPair(&TradingPair[i]))
		} else {
			json.NewEncoder(w).Encode(models.Result{"ERR", "Пара уже запущена"})
		}
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func OffPair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if i, res := SearchPair(params.Get("id")); i != -1 {
		TradingPair[i].Status = models.Off
		close(TradingPair[i].StopCh)
		json.NewEncoder(w).Encode(models.Result{"OK", "Пара остановлена"})
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func SearchPair(pairid string) (int, models.Result) {
	for i, pair := range TradingPair {
		if pairid == pair.PairId {
			return i, models.Result{Status: "OK"}
		}
	}
	return -1, models.Result{"ERR", "Пары " + pairid + " не существует"}
}
