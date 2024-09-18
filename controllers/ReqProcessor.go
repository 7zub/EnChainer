package controllers

import (
	"awesomeProject/models"
	"awesomeProject/models/exchange"
	"awesomeProject/models/exchange/exchangeReq"
	"fmt"
)

func createStruct1() []exchange.IGetReq1 {
	return []exchange.IGetReq1{
		exchange.ExchangeBinance{
			Requests: struct{ BookReq models.Request }{BookReq: exchangeReq.BinanceBookParams{}.GetParams()},
		},
	}
}

func main1() {
	structs := createStruct1()

	// Вызов метода Get для каждой структуры
	for _, obj := range structs {
		callGetMethod(obj)
	}
}

func callGetMethod(g exchange.IGetReq1) {
	fmt.Println(g.GetReq1())
}

/*
func (req models.Request) UrlCreator {
	//fields := reflect.TypeOf(req.Params)
	//values := reflect.ValueOf(req.Params)

	rq, err := http.NewRequest("GET", "https://api.binance.com/api/v3/depth", nil)
	if err != nil {
		panic(err)
	}

	//for i := 0; i < fields.NumField(); i++ {
	//	rq.URL.Query().Add(fields.Field(i).Name, "test")
	//}

	q := rq.URL.Query()
	q.Add("symbol", "SOLUSDT")
	//q.Add("param2", "value2")
	rq.URL.RawQuery = q.Encode()

	fmt.Printf("Полный URL: %s\n", rq.URL.String())
	client := http.Client{}
	resp, err := client.Do(rq)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Println(result)
}*/
