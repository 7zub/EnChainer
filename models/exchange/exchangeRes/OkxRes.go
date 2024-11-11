package exchangeRes

import (
	"awesomeProject/models"
	"strconv"
)

type OkxBook struct {
	Code string        `json:"code"`
	Msg  string        `json:"msg"`
	Data []OkxBookData `json:"data"`
}

type OkxBookData struct {
	Ts   string     `json:"ts"`
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

func (book OkxBook) Mapper() models.OrderBook {
	var newBids, newAsks []models.ValueBook

	for _, bid := range book.Data[0].Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		volume, _ := strconv.ParseFloat(bid[1], 64)

		newBids = append(newBids, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	for _, ask := range book.Data[0].Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		volume, _ := strconv.ParseFloat(ask[1], 64)

		newAsks = append(newAsks, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	return models.OrderBook{
		Exchange: models.OKX,
		Bids:     newBids,
		Asks:     newAsks,
	}
}
