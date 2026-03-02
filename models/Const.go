package models

import (
	"enchainer/controls/load"
	"time"
)

var Const Constant

type Constants struct {
	Constants Constant `yaml:"constants"`
}
type Constant struct {
	Lot            float64       `yaml:"Lot"`
	MaxTrade       int           `yaml:"MaxTrade"`
	ActiveTrade    int           `yaml:"ActiveTrade"`
	Spread         float64       `yaml:"Spread"`
	MinProfit      float64       `yaml:"MinProfit"`
	Slip           float64       `yaml:"Slip"`
	LotReserve     float64       `yaml:"LotReserve"`
	DecimalPrice   int           `yaml:"DecimalPrice"`
	DecimalVolume  int           `yaml:"DecimalVolume"`
	TimeoutBlock   time.Duration `yaml:"TimeoutBlock"`
	TimeoutCcyInfo time.Duration `yaml:"TimeoutCcyInfo"`
	BatchSize      int           `yaml:"BatchSize"`
}

func LoadConst() {
	Const = load.Yaml[Constants]("src/const.yml").Constants
}
