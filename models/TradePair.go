package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

type IResponse interface {
	Mapper() any
}

type TradePair struct {
	Id         uint `gorm:"primaryKey"`
	PairId     string
	Title      string
	Ccy        Ccy `gorm:"embedded"`
	Status     int
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

const (
	Off     = 0
	On      = 1
	Warning = 2
	Error   = 3
)

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
}

func (book OrderBook) BookExist() bool {
	if len(book.Bids) > 0 && len(book.Asks) > 0 {
		return true
	}
	return false
}

func ProfitBid(orderBooks *[]OrderBook) {

}

func ProfitAsk(orderBooks *[]OrderBook) {

}
