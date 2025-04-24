package web

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/book", BookControl)
	myRouter.HandleFunc("/addpair", AddPair)
	myRouter.HandleFunc("/deletepair", DeletePair)
	myRouter.HandleFunc("/onpair", OnPair)
	myRouter.HandleFunc("/offpair", OffPair)
	//myRouter.HandleFunc("/ws", views.Ws)
	myRouter.HandleFunc("/trade", TradeTaskControl)

	log.Fatal(http.ListenAndServe(":10", myRouter))
}
