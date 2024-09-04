package models

type TradingPair struct {
	Id        int
	Name      string
	Desc      string
	Currency  string
	Status    int
	OrderBook []OrderBook
}

type OrderBook struct {
	Exchange     int
	LastUpdateId int
	Bids         []ValueBook
	Asks         []ValueBook
}

type ValueBook struct {
	Price  float64
	Volume float64
}

const (
	Off = 0
	On  = 1
)
