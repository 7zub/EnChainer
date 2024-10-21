package models

type TradeTask struct {
	TaskId   int `gorm:"primaryKey"`
	Currency Ccy `gorm:"embedded"`
	Profit   float64
	Buy      Operation `gorm:"embedded;embeddedPrefix:buy_"`
	Sell     Operation `gorm:"embedded;embeddedPrefix:sell_"`
}

type Operation struct {
	Exchange string
	Price    float64
	Volume   *float64
}
