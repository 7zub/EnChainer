package controls

import (
	"awesomeProject/models"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"time"
)

var TradePair = []models.TradePair{}

func BookControl(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(TradePair)
}

func AddPair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	pairid := "P_" + params.Get("currency")
	sessTime, _ := strconv.ParseInt(params.Get("time"), 10, 32)
	if sessTime < 500 || sessTime > 100000 {
		json.NewEncoder(w).Encode(models.Result{"ERR", "Некорректный интервал: " + params.Get("time")})
		return
	}

	tp := models.TradePair{
		PairId: pairid,
		Title:  params.Get("title"),
		Ccy: models.Ccy{
			Currency:  params.Get("currency"),
			Currency2: "USDT",
		},
		Status:   models.Off,
		SessTime: time.Duration(sessTime) * time.Millisecond,
	}

	if i, _ := SearchPair(pairid); i == -1 {
		TradePair = append(TradePair, tp)
		SaveBookDb(&TradePair[len(TradePair)-1])

		json.NewEncoder(w).Encode(models.Result{"OK", "Пара " + params.Get("currency") + " добавлена"})
	} else {
		json.NewEncoder(w).Encode(models.Result{"ERR", "Пара " + params.Get("currency") + " уже существует"})
	}
}

func DeletePair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if i, res := SearchPair(params.Get("id")); i != -1 {
		if TradePair[i].Status == models.Off {
			DeleteBookDb(&TradePair[i])
			TradePair = slices.Delete(TradePair, i, i+1)

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
		if TradePair[i].Status != models.On {
			TradePair[i].Status = models.On
			BooksPair(&TradePair[i])
			json.NewEncoder(w).Encode(models.Result{"OK", "Мониторинг пары запущен"})
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
		TradePair[i].Status = models.Off
		close(TradePair[i].StopCh)
		json.NewEncoder(w).Encode(models.Result{"OK", "Пара остановлена"})
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func SearchPair(pairid string) (int, models.Result) {
	for i, pair := range TradePair {
		if pairid == pair.PairId {
			return i, models.Result{Status: "OK"}
		}
	}
	return -1, models.Result{"ERR", "Пары " + pairid + " не существует"}
}
