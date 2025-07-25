package controls

import (
	"enchainer/models"
	"enchainer/models/exchange"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"math"
	"math/big"
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
	data, err := os.ReadFile("rsc/config.yml")
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
	d := decimal - int(math.Floor(math.Log10(math.Abs(num)))) - 1
	pow := math.Pow10(d)
	var rounded, slip float64

	switch mode {
	case "down":
		rounded = math.Floor(num*pow) / pow
		slip = -models.Const.Split
	case "up":
		rounded = math.Ceil(num*pow) / pow
		slip = models.Const.Split
	case "near":
		rounded = math.Round(num*pow) / pow
	default:
		rounded = math.Round(num*pow) / pow
	}

	if rounded == num {
		a, _ := new(big.Float).SetString(fmt.Sprint(rounded))
		b, _ := new(big.Float).SetString(fmt.Sprint(slip / pow))
		rounded, _ = new(big.Float).Add(a, b).Float64()
	}
	return rounded
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
