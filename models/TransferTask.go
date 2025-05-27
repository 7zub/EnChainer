package models

import "time"

type TransferTask struct {
	Id     uint `gorm:"primaryKey"`
	TaskId string
	ReqId  string
	Ex     Exchange
	From   Account
	To     Account
	Ccy
	Amount     float64
	CreateDate time.Time `gorm:"type:timestamp"`
}

type Account string

const (
	Spot    Account = "Spot"
	Isolate Account = "Margin"
	Cross   Account = "Cross"
	Feature Account = "Feature"
)
