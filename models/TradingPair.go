package models

import (
	"sort"
	"time"
)

type TradingPair struct {
	Id        int
	Name      string
	Desc      string
	Ccy       Ccy
	Status    int
	SessTime  time.Duration
	OrderBook []OrderBook
}

type Ccy struct {
	Currency  string
	Currency2 string
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

func SortOrderBooks(orderBooks *[]OrderBook) {
	sort.Slice(*orderBooks, func(i, j int) bool {
		if len((*orderBooks)[i].Bids) > 0 && len((*orderBooks)[j].Bids) > 0 {
			return (*orderBooks)[i].Bids[0].Price > (*orderBooks)[j].Bids[0].Price
		} else {
			return false
		}
	})
}
