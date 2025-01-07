package models

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

type IParams[T any] interface {
	GetParams(ccy Ccy) *Request
}

type Request struct {
	ReqId       string `gorm:"primaryKey"`
	ReqType     string
	Url         string
	Params      IParams[any]   `gorm:"-"`
	Response    IResponse[any] `gorm:"-"`
	ResponseRaw string
	Code        int
	ReqDate     time.Time `gorm:"type:timestamp"`
	Log         Result    `gorm:"-"`
}

func (r *Request) SendRequest() {
	r.UrlExec(r.UrlBuild())
}

func (r *Request) DescRequest(date time.Time, rid string) {
	r.ReqDate = date
	r.ReqId = rid
	r.ReqType = "Book"
}

func GenDescRequest() (time.Time, string) {
	reqDate := time.Now()
	reqId := fmt.Sprintf("B-%02d/%02d%02d%02d-%d",
		reqDate.Day(),
		reqDate.Hour(),
		reqDate.Minute(),
		reqDate.Second(),
		rand.Intn(int(time.Now().Unix())),
	)

	return reqDate, reqId
}

func (r *Request) UrlBuild() *http.Request {
	fields := reflect.TypeOf(r.Params)
	values := reflect.ValueOf(r.Params)

	rq, err := http.NewRequest("GET", r.Url, nil)
	if err != nil {
		panic(err)
	}

	q := rq.URL.Query()

	for i := 0; i < fields.NumField(); i++ {
		q.Add(
			fields.Field(i).Tag.Get("url"),
			fmt.Sprintf("%v", values.Field(i)),
		)
	}

	rq.URL.RawQuery = q.Encode()
	return rq
}

func (r *Request) UrlExec(rq *http.Request) {
	r.Url = rq.URL.String()
	client := http.Client{}
	resp, err := client.Do(rq)
	r.Code = -1
	r.Log = Result{Status: INFO, Message: fmt.Sprintf("Запрос %s: %s", r.ReqId, rq.URL.String())}
	if err != nil {
		r.ResponseRaw = err.Error()
		r.Log = Result{Status: ERR, Message: fmt.Sprintf("Ошибка выполнения запроса %s: %s", r.ReqId, err)}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.Log = Result{Status: ERR, Message: fmt.Sprintf("Ошибка чтения ответа на %s: %s", r.ReqId, err)}
		return
	}

	err = json.Unmarshal(body, r.Response)
	if err != nil {
		r.Log = Result{Status: ERR, Message: fmt.Sprintf("Ошибка десериализации %s: %s", r.ReqId, err)}
	}
	r.ResponseRaw = string(body)
	r.Code = resp.StatusCode
}
