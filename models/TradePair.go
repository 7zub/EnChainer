package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"sort"
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
	sort.Slice(*orderBooks, func(i, j int) bool {
		if len((*orderBooks)[i].Bids) > 0 && len((*orderBooks)[j].Bids) > 0 {
			return (*orderBooks)[i].Bids[0].Price > (*orderBooks)[j].Bids[0].Price
		} else {
			return false
		}
	})

	if len(*orderBooks) > 2 {
		bestAskIndex := 0
		bestAskPrice := math.MaxFloat64

		for i, ob := range *orderBooks {
			if len(ob.Asks) > 0 && ob.Asks[0].Price < bestAskPrice {
				bestAskPrice = ob.Asks[0].Price
				bestAskIndex = i
			}
		}

		// Перемещаем элемент с лучшим asks в конец массива
		if bestAskIndex != len(*orderBooks)-1 {
			bestAskOb := (*orderBooks)[bestAskIndex]
			// Удаляем из текущей позиции и добавляем в конец
			*orderBooks = append(
				append((*orderBooks)[:bestAskIndex], (*orderBooks)[bestAskIndex+1:]...),
				bestAskOb,
			)
		}
	}
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
