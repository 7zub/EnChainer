package models

type TradeTask struct {
	TaskId   int
	Currency Ccy
	Buy      Operation
	Sell     Operation
	Profit   float64
}

type Operation struct {
	Exchange int
	Price    float64
	Volume   *float64
}
