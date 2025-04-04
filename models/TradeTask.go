package models

import "time"

type TradeTask struct {
	TaskId     string `gorm:"primaryKey"`
	Ccy        Ccy    `gorm:"embedded"`
	Spread     float64
	Buy        Operation `gorm:"embedded;embeddedPrefix:buy_"`
	Sell       Operation `gorm:"embedded;embeddedPrefix:sell_"`
	CreateDate time.Time `gorm:"type:timestamp;autoCreateTime"`
	Stage      StageTask
	Status     StatusTask
	Message    string
}

type Operation struct {
	Ex     Exchange
	Price  float64
	Volume float64
	Side   Side `gorm:"-"`
}

type OperationTask struct {
	Ccy
	Operation
}

type StageTask string

const (
	Creation   StageTask = "creation"
	Validation StageTask = "validation"
	Trade      StageTask = "trade"
)

type StatusTask string

const (
	Done     StatusTask = "done"
	Stop     StatusTask = "stop"
	Progress StatusTask = "progress"
	Err      StatusTask = "error"
)

type Side string

const (
	Buy  Side = "buy"
	Sell Side = "sell"
)
