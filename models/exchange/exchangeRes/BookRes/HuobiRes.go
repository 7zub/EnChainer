package BookRes

import (
	"enchainer/models"
)

type HuobiBook struct {
	Ch     string        `json:"ch"`
	Status string        `json:"status"`
	Ts     int64         `json:"ts"`
	Tick   HuobiBookTick `json:"tick"`
}

type HuobiBookTick struct {
	Ts      int64       `json:"ts"`
	Version int64       `json:"version"`
	Bids    [][]float64 `json:"bids"`
	Asks    [][]float64 `json:"asks"`
}

func (book HuobiBook) Mapper() any {
	var newBids, newAsks []models.ValueBook

	for _, bid := range book.Tick.Bids {
		price := bid[0]
		volume := bid[1]

		newBids = append(newBids, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	for _, ask := range book.Tick.Asks {
		price := ask[0]
		volume := ask[1]

		newAsks = append(newAsks, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	return models.OrderBook{
		Exchange: models.HUOBI,
		Bids:     newBids,
		Asks:     newAsks,
	}
}
