package models

type TradeTask struct {
	TaskId   int       `gorm:"primaryKey"`
	Currency Ccy       `gorm:"embedded"`
	Buy      Operation `gorm:"embedded"`
	Sell     Operation `gorm:"embedded"`
	Profit   float64
}

type Operation struct {
	Exchange string   `gorm:"column:exchange"`
	Price    float64  `gorm:"column:price"`
	Volume   *float64 `gorm:"column:volume"`
}
