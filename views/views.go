package views

import (
	"awesomeProject/controllers"
	"net/http"
)

func Monitor(w http.ResponseWriter, rr *http.Request) {
	controllers.Monitor(w)
}

func GetOrderBook(w http.ResponseWriter, rr *http.Request) {
	controllers.Book(w, rr)
}

func AddPair(w http.ResponseWriter, rr *http.Request) {
	controllers.AddPair(w, rr)
}

func DeletePair(w http.ResponseWriter, rr *http.Request) {
	controllers.DeletePair(w, rr)
}

func OnPair(w http.ResponseWriter, rr *http.Request) {
	controllers.OnPair(w, rr)
}

func OffPair(w http.ResponseWriter, rr *http.Request) {
	controllers.OffPair(w, rr)
}

func Ws(w http.ResponseWriter, rr *http.Request) {
	controllers.Ws()
}
