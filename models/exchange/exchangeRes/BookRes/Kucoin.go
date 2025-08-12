package BookRes

import (
	"enchainer/models"
	"strconv"
)

type KucoinBook struct {
	Code string         `json:"code"`
	Data KucoinBookData `json:"data"`
}

type KucoinBookData struct {
	Ts   int64   `json:"time"`
	Bids [][]any `json:"bids"`
	Asks [][]any `json:"asks"`
}

func (book KucoinBook) Mapper() any {
	parseLevel := func(level []any) (price, volume float64) {
		switch level[0].(type) {
		case string:
			price, _ = strconv.ParseFloat(level[0].(string), 64)
			volume, _ = strconv.ParseFloat(level[1].(string), 64)
		case float64:
			price = level[0].(float64)
			volume = level[1].(float64)
		}
		return
	}

	var newBids, newAsks []models.ValueBook

	for _, bid := range book.Data.Bids {
		p, v := parseLevel(bid)
		newBids = append(newBids, models.ValueBook{Price: p, Volume: v})
	}

	for _, ask := range book.Data.Asks {
		p, v := parseLevel(ask)
		newAsks = append(newBids, models.ValueBook{Price: p, Volume: v})
	}

	return models.OrderBook{
		Exchange: models.KUCOIN,
		Bids:     newBids,
		Asks:     newAsks,
	}
}
