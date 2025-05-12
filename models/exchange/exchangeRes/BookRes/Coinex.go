package BookRes

import (
	"enchainer/models"
	"strconv"
)

type CoinexBook struct {
	Data struct {
		Depth struct {
			Asks [][]string `json:"asks"`
			Bids [][]string `json:"bids"`
		} `json:"depth"`
	} `json:"data"`
}

func (book CoinexBook) Mapper() any {
	var newBids, newAsks []models.ValueBook

	for _, bid := range book.Data.Depth.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		volume, _ := strconv.ParseFloat(bid[1], 64)

		newBids = append(newBids, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	for _, ask := range book.Data.Depth.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		volume, _ := strconv.ParseFloat(ask[1], 64)

		newAsks = append(newAsks, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	return models.OrderBook{
		Exchange: models.COINEX,
		Bids:     newBids,
		Asks:     newAsks,
	}
}
