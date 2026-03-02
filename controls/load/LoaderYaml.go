package load

import (
	"log"
	"os"
	"path"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Loggable interface {
	GetStatus() any
	GetMessage() string
}

func Yaml[T any](filename string) T {
	data, err := os.ReadFile(filename)
	if err != nil {
		ToLog(err)
		panic("Ошибка загрузки конфигурации")
	}

	var yml T
	if err = yaml.Unmarshal(data, &yml); err != nil {
		ToLog(err)
		panic("Ошибка десериализации конфигурации")
	}
	return yml
}

func ToLog(v interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := path.Base(runtime.FuncForPC(pc).Name())

	if loggable, ok := v.(Loggable); ok {
		log.Printf("%-4s %-25s %s", loggable.GetStatus(), funcName, loggable.GetMessage())
		return
	}

	log.Printf("%-4s %-25s %v", "ERR", funcName, v)
}
