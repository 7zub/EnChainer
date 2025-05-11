package BookRes

import (
	"enchainer/models"
	"strconv"
)

type CoinexBook struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    CoinexBookData `json:"data"`
	Bids    interface{}
}

type CoinexBookData struct {
	Timestamp int64      `json:"timestamp"`
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
}

func (book CoinexBook) Mapper() any {
	var newBids, newAsks []models.ValueBook

	for _, bid := range book.Data.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		volume, _ := strconv.ParseFloat(bid[1], 64)

		newBids = append(newBids, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	for _, ask := range book.Data.Asks {
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
