package exchangeRes

import (
	"enchainer/models"
	"strconv"
)

type BybitBook struct {
	RetCode int             `json:"retCode"`
	RetMsg  string          `json:"retMsg"`
	Result  BybitBookResult `json:"result"`
}

type BybitBookResult struct {
	Bids [][]string `json:"b"`
	Asks [][]string `json:"a"`
}

func (book BybitBook) Mapper() models.OrderBook {
	var newBids, newAsks []models.ValueBook

	for _, bid := range book.Result.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		volume, _ := strconv.ParseFloat(bid[1], 64)

		newBids = append(newBids, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	for _, ask := range book.Result.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		volume, _ := strconv.ParseFloat(ask[1], 64)

		newAsks = append(newAsks, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	return models.OrderBook{
		Exchange: models.BYBIT,
		Bids:     newBids,
		Asks:     newAsks,
	}
}
