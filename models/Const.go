package models

var Const = struct {
	Lot           float64
	MaxTrade      int
	Spread        float64
	SpreadClose   float64
	Split         float64
	DecimalPrice  int
	DecimalVolume int
	BatchSize     int
}{
	Lot:           9.2,
	MaxTrade:      0,
	Spread:        0.7,
	SpreadClose:   0.1,
	Split:         3,
	DecimalPrice:  3,
	DecimalVolume: 3,
	BatchSize:     700,
}
