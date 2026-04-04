package controls

import (
	"enchainer/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

var TradeTask sync.Map
var taskIdCount atomic.Int64

func TradeTaskControl(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(TradeTask)
}

func SearchOpenTask(task *models.TradeTask) *string {
	var res *string = nil
	TradeTask.Range(func(key, val any) bool {
		t, _ := val.(*models.TradeTask)
		if task.TaskId != t.TaskId && task.Ccy == t.Ccy && t.Status != models.Stop {
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
		if ccy == t.Ccy && t.Stage == models.Trade && (t.Status == models.Pending || t.Status == models.Done) {
			res = &t.TaskId
			return false
		}
		return true
	})
	return res
}

func GenTaskId() string {
	id := taskIdCount.Add(1)
	taskId := fmt.Sprintf("T%07d", id)
	return taskId
}

func LoadTask(key string) *models.TradeTask {
	t, _ := TradeTask.Load(key)
	return t.(*models.TradeTask)
}
