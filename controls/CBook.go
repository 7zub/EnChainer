package controls

import (
	"enchainer/models"
	"encoding/json"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"time"
)

var TradePair []models.TradePair

func BookControl(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(TradePair)
}

func AddPair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	pairid := "P_" + params.Get("currency")
	sessTime, _ := strconv.ParseInt(params.Get("time"), 10, 32)
	if sessTime < 500 || sessTime > 100000 {
		json.NewEncoder(w).Encode(models.Result{Status: "ERR", Message: "Некорректный интервал: " + params.Get("time")})
		return
	}

	tp := models.TradePair{
		PairId: pairid,
		Title:  params.Get("title"),
		Market: models.Market.Spot,
		Ccy: models.Ccy{
			Currency:  params.Get("currency"),
			Currency2: "USDT",
		},
		Status:   models.Off,
		SessTime: time.Duration(sessTime+rand.Int63n(1500)) * time.Millisecond,
	}

	if i, _ := SearchPair(pairid); i == -1 {
		TradePair = append(TradePair, tp)
		SaveDb(&TradePair[len(TradePair)-1])

		json.NewEncoder(w).Encode(models.Result{Status: models.OK, Message: "Пара " + params.Get("currency") + " добавлена"})
	} else {
		json.NewEncoder(w).Encode(models.Result{Status: "ERR", Message: "Пара " + params.Get("currency") + " уже существует"})
	}
}

func DeletePair(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if i, res := SearchPair(params.Get("id")); i != -1 {
		if TradePair[i].Status == models.Off {
			DeleteBookDb(&TradePair[i])
			TradePair = slices.Delete(TradePair, i, i+1)

			json.NewEncoder(w).Encode(models.Result{Status: models.OK, Message: "Пара удалена"})
		} else {
			json.NewEncoder(w).Encode(models.Result{Status: "ERR", Message: "Пара должна быть остановлена"})
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
			SaveDb(&TradePair[i])
			StartPair(&TradePair[i])
			json.NewEncoder(w).Encode(models.Result{Status: models.OK, Message: "Мониторинг пары запущен"})
		} else {
			json.NewEncoder(w).Encode(models.Result{Status: "ERR", Message: "Пара уже запущена"})
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
		SaveDb(&TradePair[i])
		json.NewEncoder(w).Encode(models.Result{Status: models.OK, Message: "Пара остановлена"})
	} else {
		json.NewEncoder(w).Encode(res)
	}
}

func SearchPair(pairid string) (int, models.Result) {
	for i, pair := range TradePair {
		if pairid == pair.PairId {
			return i, models.Result{Status: models.OK}
		}
	}
	return -1, models.Result{Status: "ERR", Message: "Пары " + pairid + " не существует"}
}
