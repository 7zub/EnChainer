package controllers

import (
	"awesomeProject/models"
	"fmt"
	"time"
)

func TaskCreate(pair *models.TradingPair, req models.IParams) {
	ticker := time.NewTicker(pair.SessTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Выполняем запросы к биржам...")

			start := time.Now()
			r := req.GetParams(pair.Ccy)
			r.SendRequest()

			pair.OrderBook = append(pair.OrderBook, r.Response.Mapper())
			fmt.Printf("res: %v", pair.OrderBook)

			duration := time.Since(start)
			//
			//// Вывод результатов
			//for exchange, result := range taskManager.Results {
			//	fmt.Printf("Биржа: %s\n", exchange)
			//	if result.Error != nil {
			//		fmt.Printf("Ошибка: %v\n", result.Error)
			//	} else {
			//		fmt.Printf("Данные: %v\n", result.Data)
			//	}
			//	fmt.Println()
			//}

			//case <-stop:
			//	// Получен сигнал об остановке
			//	fmt.Println("Остановлен")
			//	return

			fmt.Printf("Время выполнения: %v\n. Ожидание следующего интервала...", duration)
		}
	}
}
