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
	Name       string
	Desc       string
	Ccy        Ccy `gorm:"embedded"`
	Status     int
	SessTime   time.Duration
	CreateDate time.Time
	OrderBook  []OrderBook `gorm:"foreignKey:TpID;constraint:OnDelete:CASCADE;"`
}

type Ccy struct {
	Currency  string `gorm:"column:currency"`
	Currency2 string `gorm:"column:currency2"`
}

type OrderBook struct {
	Id           uint `gorm:"primaryKey"`
	TpID         uint
	Exchange     int
	LastUpdateId int
	Bids         JsonValueBook `gorm:"type:jsonb"` //`gorm:"-"` //`gorm:"foreignKey:OrderBookID;constraint:OnDelete:CASCADE;"`
	Asks         JsonValueBook `gorm:"type:jsonb"` //`gorm:"type:jsonb"` //`gorm:"foreignKey:OrderBookID"`
}

type JsonValueBook []ValueBook

type ValueBook struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}

const (
	Off = 0
	On  = 1
)

// Scan для JSONValueBook для поддержки jsonb
func (v *JsonValueBook) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value")
	}
	return json.Unmarshal(bytes, v)
}

// Value для JSONValueBook для поддержки jsonb
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
