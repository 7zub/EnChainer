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
		d.Migrator().DropTable(
			&models.Request{},
			//&models.TradePair{},
			&models.OrderBook{},
			&models.TradeTask{},
			&models.OperationTask{},
			//&models.RequestBlock{}
			&models.TransferTask{},
		)

		err := d.AutoMigrate(
			&models.Request{},
			&models.TradePair{},
			&models.OrderBook{},
			&models.TradeTask{},
			&models.OperationTask{},
			&models.RequestBlock{},
			&models.TransferTask{})
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

func LoadBlockDb(block *[]models.RequestBlock) {
	result := db.Find(block)
	if result.Error != nil {
		ToLog(fmt.Sprintf("Ошибка БД load block: %s", result.Error))
		return
	}

	for i := range *block {
		ReqBlock.Store((*block)[i].Ccy.Currency+string((*block)[i].Ex), &(*block)[i])
	}
}
