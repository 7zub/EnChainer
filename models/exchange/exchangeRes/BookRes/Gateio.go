package BookRes

import (
	"enchainer/models"
	"strconv"
)

type GateioBook struct {
	Current float64 `json:"current"`
	Bids    []any   `json:"bids"`
	Asks    []any   `json:"asks"`
}

func (book GateioBook) Mapper() any {
	parseLevel := func(level any) (price, volume float64) {
		switch v := level.(type) {
		case []any:
			price, _ = strconv.ParseFloat(v[0].(string), 64)
			volume, _ = strconv.ParseFloat(v[1].(string), 64)
		case map[string]any:
			price, _ = strconv.ParseFloat(v["p"].(string), 64)
			volume = v["s"].(float64)
		}
		return
	}

	var bids, asks []models.ValueBook

	for _, bid := range book.Bids {
		p, v := parseLevel(bid)
		bids = append(bids, models.ValueBook{Price: p, Volume: v})
	}

	for _, ask := range book.Asks {
		p, v := parseLevel(ask)
		asks = append(asks, models.ValueBook{Price: p, Volume: v})
	}

	return models.OrderBook{
		Exchange: models.GATEIO,
		Bids:     bids,
		Asks:     asks,
	}
}
