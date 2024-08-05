package main

import (
	"awesomeProject/views"
	"github.com/gorilla/mux"
)
import "log"
import "net/http"

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", views.HomePage)
	myRouter.HandleFunc("/articles", views.ReturnAllArticles)
	myRouter.HandleFunc("/article", views.CreateNewArticle)   //.Methods("POST")
	myRouter.HandleFunc("/article/{id}", views.DeleteArticle) //.Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", views.ReturnSingleArticle)
	myRouter.HandleFunc("/hotels/", views.AllHotels)
	myRouter.HandleFunc("/hotels/export", views.ExportHotel)
	myRouter.HandleFunc("/kafka", views.Kafkatest)

	myRouter.HandleFunc("/monitor", views.Monitor)
	myRouter.HandleFunc("/book", views.GetOrderBook)
	myRouter.HandleFunc("/addpair", views.AddPair)
	myRouter.HandleFunc("/deletepair", views.DeletePair)
	myRouter.HandleFunc("/onpair", views.OnPair)
	myRouter.HandleFunc("/offpair", views.OffPair)
	log.Fatal(http.ListenAndServe(":10", myRouter))
}
