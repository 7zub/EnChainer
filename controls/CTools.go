package controls

import (
	"enchainer/models"
	"enchainer/models/exchange"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"math"
	"os"
	"path"
	"reflect"
	"runtime"
)

func ToLog(ifc interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := path.Base(runtime.FuncForPC(pc).Name())

	switch v := ifc.(type) {
	case models.Request:
		log.Printf("%-4s %-25s %s", v.Log.Status, funcName, v.Log.Message)
	case models.Result:
		log.Printf("%-4s %-25s %s", v.Status, funcName, v.Message)
	default:
		log.Printf("%-4s %-25s %s", models.ERR, funcName, v)
	}
}

func exceptTask(ex string) {
	if r := recover(); r != nil {
		ToLog(fmt.Sprintf("Паника в запросе %s: %s", ex, r))
	}
}

func LoadConf() {
	data, err := os.ReadFile("rsc/config.yaml")
	if err != nil {
		ToLog(err)
		panic("Ошибка загрузки конфигурации")
	}

	if err := yaml.Unmarshal(data, &models.Conf); err != nil {
		ToLog(err)
		panic("Ошибка десериализации конфигурации")
	}
}

func Round(num float64, decimal float64) float64 {
	factor := math.Pow(10, decimal)
	return math.Round(num*factor) / factor
}

func RoundSn(num float64, decimal int, mode string) float64 {
	order := math.Floor(math.Log10(math.Abs(num)))
	decimalPlaces := decimal - int(order) - 1
	multiplier := math.Pow10(decimalPlaces)

	switch mode {
	case "down":
		return math.Floor(num*multiplier) / multiplier
	case "up":
		return math.Ceil(num*multiplier) / multiplier
	case "near":
		return math.Round(num*multiplier) / multiplier
	default:
		return math.Round(num*multiplier) / multiplier
	}
}

func GetEx(structValue interface{}) models.Exchange {
	structType := reflect.TypeOf(structValue)
	if info, ok := exchange.ExInfo[structType]; ok {
		return info.Exchange
	}
	return "nil"
}

func GetTypeEx(exch models.Exchange, reqType string) reflect.Type {
	for typ, info := range exchange.ExInfo {
		if info.Exchange == exch && info.ReqType == reqType {
			return typ
		}
	}
	return nil
}
