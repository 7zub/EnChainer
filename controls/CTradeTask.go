package controls

import (
	"enchainer/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

var TradeTask sync.Map
var taskIdCount int

func TradeTaskControl(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(TradeTask)
}

func SearchOpenTask(task *models.TradeTask) *string {
	var res *string = nil
	TradeTask.Range(func(key, val any) bool {
		t, _ := val.(*models.TradeTask)
		if task.TaskId != t.TaskId && task.Ccy == t.Ccy && !(t.Stage == models.Trade && t.Status == models.Done) && t.Status != models.Stop {
			res = &task.TaskId
			return false
		}
		return true
	})
	return res
}

func SearchOperation(ccy models.Ccy, ex models.Exchange) (*string, int) {
	var res *string = nil
	var i int
	TradeTask.Range(func(key, val any) bool {
		t, _ := val.(*models.TradeTask)
		if ccy == t.Ccy && (ex == t.Buy.Ex || ex == t.Sell.Ex) && t.Stage == models.Trade && (t.Status == models.Pending || t.Status == models.Progress) {
			res = &t.TaskId

			if ex == t.Buy.Ex {
				i = 0
			} else if ex == t.Sell.Ex {
				i = 1
			}

			return false
		}
		return true
	})
	return res, i
}

func GenTaskId() string {
	taskIdCount = taskIdCount + 1
	taskId := fmt.Sprintf("%07d-%04d", taskIdCount, rand.Intn(10000))
	return taskId
}

func LoadTask(key string) *models.TradeTask {
	t, _ := TradeTask.Load(key)
	return t.(*models.TradeTask)
}
