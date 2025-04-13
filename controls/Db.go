package controls

import (
	"enchainer/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db = gorm.DB{}

func CreateDb() {
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

	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}})

	if err != nil {
		ToLog(err)
		panic("Не удалось подключиться к БД")
	} else {
		d.Migrator().DropTable(&models.Request{} /*&models.TradePair{},*/, &models.OrderBook{}, &models.TradeTask{}, &models.OperationTask{}, &models.RequestBlock{})

		err := d.AutoMigrate(&models.Request{}, &models.TradePair{}, &models.OrderBook{}, &models.TradeTask{}, &models.OperationTask{}, &models.RequestBlock{})
		if err != nil {
			ToLog(err)
			panic("Ошибка миграции БД")
		}
		db = *d
	}
}

func SaveDb(obj any) {
	result := db.Save(obj)

	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД %T: %s", obj, result.Error))
	}
}

func SaveBookDb(pair *models.TradePair) {
	result := db.Save(pair)

	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД order book: %s, %s", result.Error, pair.PairId))
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

func SaveTradeTaskDb(task *models.TradeTask) {
	result := db.Save(&task)

	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД task: %s", result.Error))
	}
}

func LoadTradeTaskDb(task *models.TradeTask) {
	result := db.Find(task)
	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД load task: %s", result.Error))
	}
}

func SaveReqDb(req *models.Request) {
	result := db.Save(&req)

	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД request: %s", result.Error))
	}
}

func SaveReqBlockDb(reqb *models.RequestBlock) {
	result := db.Save(&reqb)

	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД RequestBlock: %s", result.Error))
	}
}
