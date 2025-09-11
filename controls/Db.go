package controls

import (
	"enchainer/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var db = gorm.DB{}
var ChanBook = make(chan []models.OrderBook, 1000)
var ChanAny = make(chan any, 100)

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
			//&models.RequestBlock{}
			//&models.TradePair{},
			&models.OrderBook{},
			&models.TradeTask{},
			&models.OperationTask{},
			&models.TransferTask{},
		)

		err := d.AutoMigrate(
			&models.Request{},
			&models.RequestBlock{},
			&models.TradePair{},
			&models.OrderBook{},
			&models.TradeTask{},
			&models.OperationTask{},
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
		ReqBlock.Store(string((*block)[i].Market)+(*block)[i].Ccy.Currency+string((*block)[i].Ex), &(*block)[i])
	}
}

func DbSaver(ch1 <-chan []models.OrderBook, ch2 <-chan any) {
	batch := make([]models.OrderBook, 0, 1000)
	ticker := time.NewTicker(time.Second * 200)
	defer ticker.Stop()

	for {
		select {
		case p := <-ch2:
			result := db.Save(p)

			if result.Error != nil {
				ToLog(fmt.Sprintf("Ошибка БД %T: %s", p, result.Error))
			}

		case ob := <-ch1:
			batch = append(batch, ob...)
			if len(batch) >= models.Const.BatchSize {
				result := db.Save(batch)
				if result.Error != nil {
					ToLog(fmt.Sprintf("Ошибка БД при сохранении batch %T: %s", batch, result.Error))
				} else {
					ToLog(models.Result{Status: models.INFO, Message: fmt.Sprintf("Сохранен batch %T: размером %v", batch, len(batch))})
				}
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				result := db.Save(&batch)
				batch = batch[:0]

				if result.Error != nil {
					ToLog(fmt.Sprintf("Ошибка БД при сохранении batch %T: %s", batch, result.Error))
				}
			}
		}
	}
}
