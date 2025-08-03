package models

var Const = struct {
	Lot         float64
	MaxTrade    int
	Spread      float64
	SpreadClose float64
	Split       float64
	BatchSize   int
}{
	Lot:         9.2,
	MaxTrade:    2,
	Spread:      0.7,
	SpreadClose: 0.1,
	Split:       3,
	BatchSize:   300,
}
