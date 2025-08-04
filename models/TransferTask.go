package models

import "time"

type TransferTask struct {
	Id     uint `gorm:"primaryKey"`
	TaskId string
	ReqId  string
	Ex     Exchange
	From   MarketType
	To     MarketType
	Ccy
	Amount     float64
	CreateDate time.Time `gorm:"type:timestamp"`
}
