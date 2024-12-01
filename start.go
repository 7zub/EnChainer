// main.go
package main

import (
	"enchainer/controls"
	"log"
	"os"
)

func main() {
	logOn()
	controls.CreateDb()
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
