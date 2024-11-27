package controls

import (
	"enchainer/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db = gorm.DB{}

func CreateDb() {
	dsn := "host=localhost user=postgres password=Lost4096## dbname=postgres port=5432 search_path=ex sslmode=disable"
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к БД")
	} else {
		d.Migrator().DropTable(&models.Request{}, &models.TradePair{}, &models.OrderBook{}, &models.TradeTask{})

		err := d.AutoMigrate(&models.Request{}, &models.TradePair{}, &models.OrderBook{}, &models.TradeTask{})
		if err != nil {
			panic("Ошибка миграции БД")
		}
		db = *d
	}
}

func SaveBookDb(pair *models.TradePair) {
	result := db.Save(&pair)

	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД order book: %s", result.Error))
	}
}

func DeleteBookDb(pair *models.TradePair) {
	result := db.Delete(&pair)

	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД order book: %s", result.Error))
	}
}

func SaveTradeDb(task *models.TradeTask) {
	result := db.Save(&task)

	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД task: %s", result.Error))
	}
}

func SaveReqDb(req *models.Request) {
	result := db.Save(&req)

	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД request: %s", result.Error))
	}
}
