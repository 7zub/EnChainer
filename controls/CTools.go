package controls

import (
	"awesomeProject/models"
	"fmt"
	"log"
	"runtime"
)

type Any struct {
	value interface{}
}

func (v Any) Str() string {
	return fmt.Sprintf("%v", v.value)
}

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

func exceptTask() {
	if r := recover(); r != nil {
		ToLog(fmt.Sprintf("Паника: %s", r))
	}
}
