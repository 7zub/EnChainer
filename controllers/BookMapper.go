package controllers

import (
	"awesomeProject/models"
	"awesomeProject/models/exchangeResponse"
	"strconv"
)

func BookMapper(book exchangeResponse.BinanceBook) models.OrderBook {
	var ll []models.ValueBook

	for _, bid := range book.Bids {
		price, _ := strconv.ParseFloat(bid[0], 64)
		volume, _ := strconv.ParseFloat(bid[1], 64)

		ll = append(ll, models.ValueBook{
			Price:  price,
			Volume: volume,
		})
	}

	return models.OrderBook{
		LastUpdateId: book.LastUpdateId,
		Bids:         ll,
	}
}
