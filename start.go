// main.go
package main

import (
	_ "errors"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("\a")
	handleRequests()

	//filename, err := os.Open("hello.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer filename.Close()
	//
	//data, err := io.ReadAll(filename)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//var result []models.Article
	//
	//jsonErr := json.Unmarshal(data, &result)
	//
	//if jsonErr != nil {
	//	log.Fatal(jsonErr)
	//}
	//
	//fmt.Println(result)
}
