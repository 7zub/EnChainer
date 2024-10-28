package models

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type IParams interface {
	GetParams(ccy Ccy) *Request
}

type Request struct {
	ReqId    string `gorm:"primaryKey"`
	Url      string
	Params   IParams   `gorm:"-"`
	Response IResponse `gorm:"-"`
	ReqDate  time.Time `gorm:"type:timestamp"`
}

func (r *Request) SendRequest() {
	r.UrlExec(r.UrlBuild())
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
			strings.ToLower(fields.Field(i).Name),
			fmt.Sprintf("%v", values.Field(i)),
		)
	}

	rq.URL.RawQuery = q.Encode()
	log.Printf("Полный URL: %s\n", rq.URL.String())
	fmt.Printf("Полный URL: %s\n", rq.URL.String())
	return rq
}

func (r *Request) UrlExec(rq *http.Request) {
	r.ReqDate = time.Now()
	r.Url = rq.URL.String()
	r.ReqId = fmt.Sprintf("B-%02d%02d%02d%02d%03d%03d",
		r.ReqDate.Day(),
		r.ReqDate.Hour(),
		r.ReqDate.Minute(),
		r.ReqDate.Second(),
		r.ReqDate.Nanosecond()/1e6,
		rand.Intn(1000),
	)
	client := http.Client{}
	resp, err := client.Do(rq)
	if err != nil {
		log.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(r.Response)
	if err != nil {
		log.Println(err)
		return
	}
}
