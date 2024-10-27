package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type TradingPair struct {
	Id         uint `gorm:"primaryKey"`
	PairId     string
	Title      string
	Ccy        Ccy `gorm:"embedded"`
	Status     int
	SessTime   time.Duration
	CreateDate time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdateDate time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	OrderBook  []OrderBook   `gorm:"foreignKey:TpId;constraint:OnDelete:CASCADE;"`
	StopCh     chan struct{} `gorm:"-"`
}

type Ccy struct {
	Currency  string `gorm:"column:currency"`
	Currency2 string `gorm:"column:currency2"`
}

type OrderBook struct {
	Id       uint `gorm:"primaryKey"`
	TpId     uint
	Exchange string
	Bids     JsonValueBook `gorm:"type:jsonb"`
	Asks     JsonValueBook `gorm:"type:jsonb"`
	ReqDate  time.Time     `gorm:"type:timestamp"`
	ResDate  time.Time     `gorm:"type:timestamp;autoCreateTime"`
	ReqId    string
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
