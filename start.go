package main

import (
	"enchainer/controls"
	"enchainer/models"
	"log"
	"os"
)

func main() {
	logOn()
	controls.CreateDb()
	controls.LoadBookDb(&controls.TradePair)
	for i, pair := range controls.TradePair {
		if pair.Status == models.On {
			controls.BooksPair(&controls.TradePair[i])
		}
	}
	handleRequests()
}

func logOn() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(file)
}
