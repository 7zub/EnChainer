package models

type TradeTask struct {
	TaskId       int
	Currency     string
	ExchangeBuy  int
	ExchangeSell int
	PriceBuy     float64
	PriceSell    float64
	Profit       float64
}
