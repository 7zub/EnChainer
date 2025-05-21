package models

import (
	"sync"
	"time"
)

type TradeTask struct {
	TaskId     string `gorm:"primaryKey"`
	Ccy        Ccy    `gorm:"embedded"`
	Spread     float64
	Buy        Operation       `gorm:"embedded;embeddedPrefix:buy_"`
	Sell       Operation       `gorm:"embedded;embeddedPrefix:sell_"`
	OpTask     []OperationTask `gorm:"foreignKey:TaskId;constraint:OnDelete:CASCADE;"`
	CreateDate time.Time       `gorm:"type:timestamp;autoCreateTime"`
	Stage      StageTask
	Status     StatusTask
	Message    string
	Mu         sync.RWMutex `gorm:"-"`
}

type StageTask string
type StatusTask string

const (
	Creation   StageTask = "creation"
	Validation StageTask = "validation"
	Trade      StageTask = "trade"
)

const (
	Done     StatusTask = "done"
	Stop     StatusTask = "stop"
	Pending  StatusTask = "pending"
	Progress StatusTask = "progress"
	Err      StatusTask = "error"
)
