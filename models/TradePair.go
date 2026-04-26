package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"
)

type IResponse interface {
	Mapper() any
}

type TradePair struct {
	Id         uint `gorm:"primaryKey;autoIncrement"`
	PairId     string
	Title      string
	Market     MarketType
	Ccy        Ccy `gorm:"embedded"`
	Status     StPair
	SessTime   time.Duration
	CreateDate time.Time      `gorm:"type:timestamp;autoCreateTime"`
	UpdateDate time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	OrderBook  []OrderBook    `gorm:"foreignKey:TpId;constraint:OnDelete:CASCADE;"`
	StopCh     chan struct{}  `gorm:"-"`
	Mu         sync.Mutex     `gorm:"-"`
	Wg         sync.WaitGroup `gorm:"-"`
}

type Ccy struct {
	Currency  string `gorm:"column:ccy"`
	Currency2 string `gorm:"column:ccy2"`
}

type MarketType string

var Market = struct {
	Spot, Isolate, Cross, Futures MarketType
}{
	Spot:    "spot",
	Isolate: "margin",
	Cross:   "cross",
	Futures: "futures",
}

type OrderBook struct {
	Id         uint `gorm:"primaryKey"`
	TpId       uint
	Exchange   Exchange
	Bids       JsonValueBook `gorm:"type:jsonb"`
	Asks       JsonValueBook `gorm:"type:jsonb"`
	CreateDate time.Time     `gorm:"type:timestamp;autoCreateTime"`
	ReqId      string        `gorm:"unique"`
	TaskId     string
}

type JsonValueBook []ValueBook

type ValueBook struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}

type StPair string

var StatusPair = struct {
	On, Off, Pause StPair
}{
	On: "On", Off: "Off", Pause: "Pause",
}

func (v *JsonValueBook) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value")
	}
	return json.Unmarshal(bytes, v)
}

func (v JsonValueBook) Value() (driver.Value, error) {
	return json.Marshal(v)
}

func SortOrderBooks(orderBooks *[]OrderBook) {
	obs := *orderBooks
	bestBidIdx, bestAskIdx := 0, 0
	bestBid := -1.0
	bestAsk := math.MaxFloat64

	for i, ob := range obs {
		if len(ob.Bids) > 0 && ob.Bids[0].Price > bestBid {
			bestBid = ob.Bids[0].Price
			bestBidIdx = i
		}
		if len(ob.Asks) > 0 && ob.Asks[0].Price < bestAsk {
			bestAsk = ob.Asks[0].Price
			bestAskIdx = i
		}
	}

	obs[0], obs[bestBidIdx] = obs[bestBidIdx], obs[0]

	if bestAskIdx == bestBidIdx {
		bestAskIdx = 0
	} else if bestAskIdx == 0 {
		bestAskIdx = bestBidIdx
	}

	if bestAskIdx != 0 {
		last := len(obs) - 1
		obs[last], obs[bestAskIdx] = obs[bestAskIdx], obs[last]
	}

	*orderBooks = obs
}

func (book OrderBook) BookExist() bool {
	if len(book.Bids) > 0 && len(book.Asks) > 0 {
		return true
	}
	return false
}

func GetVolume(valueBook *JsonValueBook) (ValueBook, int) {
	var p, v, usd float64
	var deep int

	for i, book := range *valueBook {
		p += book.Price
		v += book.Volume
		usd += book.Price * book.Volume
		deep = i + 1

		if usd > Const.Lot*Const.LotReserve {
			return ValueBook{Price: p / float64(deep), Volume: v}, deep
		}
	}

	if deep == 0 {
		deep = 1
	}
	return ValueBook{Price: p / float64(deep), Volume: v}, deep
}
