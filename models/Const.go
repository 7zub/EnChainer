package models

import "time"

var Const = struct {
	Lot           float64
	MaxTrade      int
	Spread        float64
	SpreadClose   float64
	Split         float64
	DecimalPrice  int
	DecimalVolume int
	TimeoutBlock  time.Duration
	BatchSize     int
}{
	Lot:           10.1,
	MaxTrade:      1,
	Spread:        0.4,
	SpreadClose:   0.1,
	Split:         3,
	DecimalPrice:  4,
	DecimalVolume: 2,
	TimeoutBlock:  160,
	BatchSize:     700,
}
