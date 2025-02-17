package web

import (
	"enchainer/_dev"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", _dev.HomePage)
	myRouter.HandleFunc("/articles", _dev.ReturnAllArticles)
	myRouter.HandleFunc("/article", _dev.CreateNewArticle)   //.Methods("POST")
	myRouter.HandleFunc("/article/{id}", _dev.DeleteArticle) //.Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", _dev.ReturnSingleArticle)
	myRouter.HandleFunc("/hotels/", _dev.AllHotels)
	myRouter.HandleFunc("/hotels/export", _dev.ExportHotel)
	myRouter.HandleFunc("/kafka", _dev.Kafkatest)

	myRouter.HandleFunc("/book", BookControl)
	myRouter.HandleFunc("/addpair", AddPair)
	myRouter.HandleFunc("/deletepair", DeletePair)
	myRouter.HandleFunc("/onpair", OnPair)
	myRouter.HandleFunc("/offpair", OffPair)
	//myRouter.HandleFunc("/ws", views.Ws)
	myRouter.HandleFunc("/trade", TradeTaskControl)

	log.Fatal(http.ListenAndServe(":10", myRouter))
}
