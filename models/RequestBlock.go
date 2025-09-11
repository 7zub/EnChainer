package models

import (
	"time"
)

type RequestBlock struct {
	Id         int `gorm:"primaryKey"`
	ReqId      string
	Market     MarketType
	Ccy        Ccy `gorm:"embedded"`
	Ex         Exchange
	ReasonCode int
	Reason     string
	CreateDate time.Time `gorm:"type:timestamp;autoCreateTime"`
	RepeatDate time.Time `gorm:"type:timestamp"`
	Active     bool
}
