package models

import "time"

type Operation struct {
	Ex     Exchange
	Price  float64
	Volume float64
	Side   Side
}

type OperationTask struct {
	Id     uint `gorm:"primaryKey"`
	TaskId string
	ReqId  string
	Ccy
	Operation
	Commission float32
	CreateDate time.Time `gorm:"type:timestamp"`
}

type Side string

const (
	Buy  Side = "buy"
	Sell Side = "sell"
)
