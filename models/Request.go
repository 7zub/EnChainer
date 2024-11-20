package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

type IParams interface {
	GetParams(ccy Ccy) *Request
}

type Request struct {
	ReqId       string `gorm:"primaryKey"`
	Url         string
	Params      IParams   `gorm:"-"`
	Response    IResponse `gorm:"-"`
	ResponseRaw string
	Code        int
	ReqDate     time.Time `gorm:"type:timestamp"`
	Log         Result    `gorm:"-"`
}

func (r *Request) SendRequest() {
	r.UrlExec(r.UrlBuild())
}

func (r *Request) DescRequest() {
	r.ReqDate = time.Now()
	r.ReqId = fmt.Sprintf("B-%02d%02d%02d%02d%03d%03d",
		r.ReqDate.Day(),
		r.ReqDate.Hour(),
		r.ReqDate.Minute(),
		r.ReqDate.Second(),
		r.ReqDate.Nanosecond()/1e6,
		rand.Intn(1000),
	)
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
	r.Log = Result{Status: INFO, Message: fmt.Sprintf("Полный URL: %s\n", rq.URL.String())}
	return rq
}

func (r *Request) UrlExec(rq *http.Request) {
	r.Url = rq.URL.String()
	client := http.Client{}
	resp, err := client.Do(rq)
	r.Code = -1

	if err != nil {
		r.ResponseRaw = err.Error()
		log.Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(body, r.Response)
	if err != nil {
		log.Println(err)
	}
	r.ResponseRaw = string(body)
	r.Code = resp.StatusCode
}
