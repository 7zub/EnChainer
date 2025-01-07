package models

import "time"

type TradeTask struct {
	TaskId     int `gorm:"primaryKey"`
	Ccy        Ccy `gorm:"embedded"`
	Spread     float64
	Buy        Operation `gorm:"embedded;embeddedPrefix:buy_"`
	Sell       Operation `gorm:"embedded;embeddedPrefix:sell_"`
	CreateDate time.Time `gorm:"type:timestamp;autoCreateTime"`
}

type Operation struct {
	Ex     string
	Price  float64
	Volume float64
}
