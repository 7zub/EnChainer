package exchangeRes

import (
	"enchainer/models"
	"strconv"
)

type KucoinBook struct {
	Code string         `json:"code"`
	Data KucoinBookData `json:"data"`
}

type KucoinBookData struct {
	Ts   int64      `json:"time"`
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

func (book KucoinBook) Mapper() models.OrderBook {
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
		Exchange: models.KUCOIN,
		Bids:     newBids,
		Asks:     newAsks,
	}
}
