package models

import "time"

var Const = struct {
	Lot           float64
	MaxTrade      int
	Spread        float64
	MinProfit     float64
	Slip          float64
	DecimalPrice  int
	DecimalVolume int
	TimeoutBlock  time.Duration
	BatchSize     int
}{
	Lot:           24.1,
	MaxTrade:      2,
	Spread:        0.49,
	MinProfit:     0.3,
	Slip:          3,
	DecimalPrice:  4,
	DecimalVolume: 2,
	TimeoutBlock:  160,
	BatchSize:     700,
}
