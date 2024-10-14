package controls

import (
	"awesomeProject/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db = gorm.DB{}

func CreateDb() {
	dsn := "host=localhost user=postgres password=Lost4096## dbname=postgres port=5432 search_path=ex sslmode=disable"
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	} else {
		d.Migrator().DropTable(&models.TradingPair{}, &models.OrderBook{}, &models.TradeTask{})

		err := d.AutoMigrate(&models.TradingPair{}, &models.OrderBook{}, &models.TradeTask{})
		if err != nil {
			panic("failed to migrate database")
		}
		db = *d
	}
}

func SaveBookDb(pair *models.TradingPair) {
	result := db.Save(&pair)

	if result.Error != nil {
		fmt.Println("Error creating order book:", result.Error)
	} else {
		fmt.Printf("OrderBook created successfully with ID: %d\n", pair.Ccy)
	}
}

func DeleteBookDb(pair *models.TradingPair) {
	result := db.Delete(&pair)

	if result.Error != nil {
		fmt.Println("Error creating order book:", result.Error)
	} else {
		fmt.Printf("OrderBook created successfully with ID: %d\n", pair.Ccy)
	}
}

func SaveTradeDb(task *models.TradeTask) {
	result := db.Save(&task)

	if result.Error != nil {
		fmt.Println("Error creating task:", result.Error)
	} else {
		fmt.Printf("Task created successfully with ID: %d\n", task.TaskId)
	}
}
