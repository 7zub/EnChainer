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
