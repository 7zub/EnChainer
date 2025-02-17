package controls

import (
	"enchainer/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db = gorm.DB{}

func CreateDb() {
	//dsn := "host={localhost} user=postgres password=Lost4096## dbname=postgres port=5432 search_path=ex sslmode=disable"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d search_path=%s sslmode=%s",
		models.Conf.Db.Host,
		models.Conf.Db.User,
		models.Conf.Db.Password,
		models.Conf.Db.Name,
		models.Conf.Db.Port,
		models.Conf.Db.Path,
		models.Conf.Db.SslMode,
	)

	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		ToLog(err)
		panic("Не удалось подключиться к БД")
	} else {
		d.Migrator().DropTable(&models.Request{} /*&models.TradePair{},*/, &models.OrderBook{}, &models.TradeTask{})

		err := d.AutoMigrate(&models.Request{}, &models.TradePair{}, &models.OrderBook{}, &models.TradeTask{})
		if err != nil {
			ToLog(err)
			panic("Ошибка миграции БД")
		}
		db = *d
	}
}

func SaveBookDb(pair *models.TradePair) {
	result := db.Save(pair)

	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД order book: %s", result.Error))
	}
}

func LoadBookDb(pairs *[]models.TradePair) {
	result := db.Find(pairs)
	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД load book: %s", result.Error))
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
