// main.go
package main

import (
	"awesomeProject/controls"
	"log"
	"os"
)

func main() {
	logOn()
	controls.CreateDb()
	handleRequests()
}

func logOn() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)
}
