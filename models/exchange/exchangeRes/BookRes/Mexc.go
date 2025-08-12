package BookRes

import (
	"enchainer/models"
	"strconv"
)

type MexcBook struct {
	Code int          `json:"code"`
	Bids [][]any      `json:"bids"`
	Asks [][]any      `json:"asks"`
	Data MexcBookData `json:"data"`
}

type MexcBookData struct {
	Bids [][]any `json:"bids"`
	Asks [][]any `json:"asks"`
}

func (book MexcBook) Mapper() any {
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
	var bids, asks [][]any
	if len(book.Data.Bids) > 0 {
		bids = book.Data.Bids
		asks = book.Data.Asks
	} else {
		bids = book.Bids
		asks = book.Asks
	}
	for _, bid := range bids {
		p, v := parseLevel(bid)
		newBids = append(newBids, models.ValueBook{Price: p, Volume: v})
	}

	for _, ask := range asks {
		p, v := parseLevel(ask)
		newAsks = append(newBids, models.ValueBook{Price: p, Volume: v})
	}

	return models.OrderBook{
		Exchange: models.MEXC,
		Bids:     newBids,
		Asks:     newAsks,
	}
}
