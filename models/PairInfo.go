package models

import (
	"time"
)

type PairInfo struct {
	Id         uint `gorm:"primaryKey;autoIncrement"`
	Ex         Exchange
	Ccy        Ccy `gorm:"embedded"`
	Cct        float64
	Lever      int
	CreateDate time.Time  `gorm:"type:timestamp;autoCreateTime"`
	UpdateDate *time.Time `gorm:"type:timestamp"`
	ReloadDate time.Time  `gorm:"type:timestamp;autoUpdateTime"`
}
