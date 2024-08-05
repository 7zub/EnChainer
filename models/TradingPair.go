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
	LastUpdateId int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

const (
	Off = 0
	On  = 1
)

const (
	BINANCE = 1
	GATEIO  = 2
	BYBIT
)
