package web

import (
	"enchainer/controls"
	"net/http"
)

func BookControl(w http.ResponseWriter, rr *http.Request) {
	controls.BookControl(w)
}

func AddPair(w http.ResponseWriter, rr *http.Request) {
	controls.AddPair(w, rr)
}

func DeletePair(w http.ResponseWriter, rr *http.Request) {
	controls.DeletePair(w, rr)
}

func OnPair(w http.ResponseWriter, rr *http.Request) {
	controls.OnPair(w, rr)
}

func OffPair(w http.ResponseWriter, rr *http.Request) {
	controls.OffPair(w, rr)
}

func Ws(w http.ResponseWriter, rr *http.Request) {
	controls.Ws()
}

func TradeTaskControl(w http.ResponseWriter, rr *http.Request) {
	controls.TradeTaskControl(w)
}

func Settings(w http.ResponseWriter, rr *http.Request) {
	controls.Settings(w, rr)
}
