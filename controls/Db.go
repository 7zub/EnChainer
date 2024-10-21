package controls

import (
	"awesomeProject/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db = gorm.DB{}

func CreateDb() {
	dsn := "host=localhost user=postgres password=Lost4096## dbname=postgres port=5432 search_path=ex sslmode=disable"
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	} else {
		d.Migrator().DropTable(&models.Request{}, &models.TradingPair{}, &models.OrderBook{}, &models.TradeTask{})

		err := d.AutoMigrate(&models.Request{}, &models.TradingPair{}, &models.OrderBook{}, &models.TradeTask{})
		if err != nil {
			panic("failed to migrate database")
		}
		db = *d
	}
}

func SaveBookDb(pair *models.TradingPair) {
	result := db.Save(&pair)

	if result.Error != nil {
		log.Println("Error creating order book:", result.Error)
	} else {
		log.Printf("OrderBook created successfully with ID: %d\n", pair.Ccy.Currency)
	}
}

func DeleteBookDb(pair *models.TradingPair) {
	result := db.Delete(&pair)

	if result.Error != nil {
		log.Println("Error creating order book:", result.Error)
	} else {
		log.Printf("OrderBook created successfully with ID: %d\n", pair.Ccy)
	}
}

func SaveTradeDb(task *models.TradeTask) {
	result := db.Save(&task)

	if result.Error != nil {
		log.Println("Error creating task:", result.Error)
	} else {
		log.Printf("Task created successfully with ID: %d\n", task.TaskId)
	}
}

func SaveReqDb(req *models.Request) {
	result := db.Save(&req)

	if result.Error != nil {
		log.Println("Error creating request:", result.Error)
	} else {
		log.Printf("Request created successfully with ID: %d\n", req.ReqId)
	}
}
