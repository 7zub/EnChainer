package controls

import (
	"awesomeProject/models"
	"fmt"
	"log"
	"runtime"
)

func ToLog(ifc interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	switch v := ifc.(type) {
	case models.Request:
		log.Printf("%s %s, %s", v.Log.Status, funcName, v.Log.Message)
	default:
		log.Printf("%s, %s, %s", models.ERR, funcName, v)
	}
}

func exceptTask(ex string) {
	if r := recover(); r != nil {
		ToLog(fmt.Sprintf("Паника в запросе %s: %s", ex, r))
	}
}
