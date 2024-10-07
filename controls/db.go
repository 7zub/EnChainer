package controls

import (
	"awesomeProject/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// OrderBook представляет книгу ордеров с полями для Bids и Asks
type OrderBook struct {
	ID           uint `gorm:"primaryKey"` // Уникальный идентификатор
	Exchange     int  // Биржа
	LastUpdateId int  // Последний идентификатор обновления
	Tgr          string
	Bids         []ValueBook `gorm:"foreignKey:OrderBookID;constraint:OnDelete:CASCADE;"` // Bid ордера
	Asks         []ValueBook `gorm:"foreignKey:OrderBookID;constraint:OnDelete:CASCADE;"` // Ask ордера
}

// ValueBook представляет ценовые значения и объемы (Bid или Ask)
type ValueBook struct {
	ID          uint    `gorm:"primaryKey"` // Уникальный идентификатор
	OrderBookID uint    // Внешний ключ на таблицу OrderBook
	Price       float64 // Цена ордера
	Volume      float64 // Объем ордера
}

func db() {
	// Подключение к базе данных PostgreSQL
	dsn := "host=localhost user=postgres password= dbname=postgres port=5432 search_path=ex sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Миграция схемы базы данных (создание таблиц, если их ещё нет)
	err = db.AutoMigrate(&OrderBook{}, &ValueBook{})
	if err != nil {
		panic("failed to migrate database")
	}

	// Пример создания новой записи OrderBook с Bids и Asks
	orderBook := OrderBook{
		Exchange:     1,
		LastUpdateId: 12345,
		Bids: []ValueBook{
			{Price: 100.5, Volume: 10},
			{Price: 101.0, Volume: 5},
		},
		Asks: []ValueBook{
			{Price: 102.0, Volume: 8},
			{Price: 103.5, Volume: 3},
		},
	}

	// Сохранение записи в базу данных
	result := db.Create(&orderBook)

	// Проверка успешности сохранения
	if result.Error != nil {
		fmt.Println("Error creating order book:", result.Error)
	} else {
		fmt.Printf("OrderBook created successfully with ID: %d\n", orderBook.ID)
	}
}

type Base struct {
	ID    uint `gorm:"primaryKey"`
	Table struct{}
}

func db1(ppp *models.TradingPair) {
	dsn := "host=localhost user=postgres password= dbname=postgres port=5432 search_path=ex sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	//b := Base{
	//	ID:    0,
	//	Table: &ppp,
	//}

	// Миграция схемы базы данных (создание таблиц, если их ещё нет)
	//err = db.AutoMigrate(&Base{})
	err = db.AutoMigrate(&models.TradingPair{}, &models.OrderBook{}) //, &models.ValueBook{})
	if err != nil {
		panic("failed to migrate database")
	}

	// Сохранение записи в базу данных
	result := db.Save(&ppp)

	// Проверка успешности сохранения
	if result.Error != nil {
		fmt.Println("Error creating order book:", result.Error)
	} else {
		fmt.Printf("OrderBook created successfully with ID: %d\n", ppp.Ccy)
	}

}
