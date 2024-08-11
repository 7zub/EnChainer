package controllers

import (
	"awesomeProject/models"
	"awesomeProject/models/exchangeResponse"
	"strconv"
)

func BookMapper(book exchangeResponse.BinanceBook) models.OrderBook {
	var newBids, newAsks []models.ValueBook

	for _, bid := range book.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		volume, _ := strconv.ParseFloat(bid[1], 64)

		newBids = append(newBids, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	for _, ask := range book.Asks {
		price, _ := strconv.ParseFloat(ask[0], 64)
		volume, _ := strconv.ParseFloat(ask[1], 64)

		newAsks = append(newAsks, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	return models.OrderBook{
		Exchange:     models.BINANCE,
		LastUpdateId: book.LastUpdateId,
		Bids:         newBids,
		Asks:         newAsks,
	}
}
