package main

import (
	"enchainer/controls"
	"enchainer/models"
	"enchainer/web"
	"log"
	"os"
	"time"
)

func main() {
	logOn()
	controls.LoadConf()
	controls.CreateDb()
	controls.LoadBlockDb(&[]models.RequestBlock{})
	controls.LoadBookDb(&controls.TradePair)

	for i, pair := range controls.TradePair {
		if pair.Status == models.On {
			controls.StartPair(&controls.TradePair[i])
			time.Sleep(300 * time.Millisecond)
		}
	}
	//controls.Trade()
	web.HandleRequests()
}

func logOn() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(file)
}
