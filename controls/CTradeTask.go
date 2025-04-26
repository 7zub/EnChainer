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

func SearchPendTask(ccy models.Ccy) *string {
	var res *string = nil
	TradeTask.Range(func(key, val any) bool {
		t, _ := val.(*models.TradeTask)
		if ccy == t.Ccy && t.Stage == models.Trade && (t.Status == models.Pending || t.Status == models.Progress) {
			res = &t.TaskId
			return false
		}
		return true
	})
	return res
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
