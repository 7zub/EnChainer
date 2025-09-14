package models

import "time"

var Const = struct {
	Lot           float64
	MaxTrade      int
	Spread        float64
	MinProfit     float64
	Slip          float64
	LotReserve    float64
	DecimalPrice  int
	DecimalVolume int
	TimeoutBlock  time.Duration
	BatchSize     int
}{
	Lot:           19.1,
	MaxTrade:      0,
	Spread:        0.65,
	MinProfit:     0.3,
	Slip:          3,
	LotReserve:    1.1,
	DecimalPrice:  4,
	DecimalVolume: 2,
	TimeoutBlock:  160,
	BatchSize:     700,
}
