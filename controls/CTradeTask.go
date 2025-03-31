package controls

import (
	"enchainer/models"
	"encoding/json"
	"net/http"
)

var TradeTask []models.TradeTask

func TradeTaskControl(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(TradeTask)
}

func SearchOpenTask(task models.TradeTask) int {
	for i, t := range TradeTask {
		if task.TaskId != t.TaskId && task.Ccy == t.Ccy && !(t.Stage == models.Trade && t.Status == models.Done) && t.Status != models.Stop {
			return i
		}
	}
	return -1
}
