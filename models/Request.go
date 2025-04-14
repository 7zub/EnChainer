package models

import (
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

var reqIdCount int

type IParams interface {
	GetParams(t any) *Request
}

type Request struct {
	ReqId       string `gorm:"primaryKey"`
	ReqType     string
	SignWay     func(r *http.Request) `gorm:"-"`
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

func (r *Request) DescRequest(date time.Time, rid string) {
	r.ReqDate = date
	r.ReqId = rid
}

func GenDescRequest() (time.Time, string) {
	reqDate := time.Now()
	reqIdCount = reqIdCount + 1
	reqId := fmt.Sprintf("%07d-%04d", reqIdCount, rand.Intn(10000))
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
		tag := fields.Field(i).Tag.Get("url")

		if tag != "-" {
			q.Add(tag, fmt.Sprintf("%v", values.Field(i)))
		}
	}
	rq.URL.RawQuery = q.Encode()

	switch r.ReqType {
	case "Trade", "Balance":
		rq.Method = "POST"
		r.SignWay(rq)
	}

	rq.URL.RawQuery = rq.URL.Query().Encode()
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

	r.Code = resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	r.ResponseRaw = string(body)

	if resp.StatusCode == 429 {
		r.Log = Result{Status: WAR, Message: fmt.Sprintf("Превышен лимит запросов к api %s", r.ReqId)}
		return
	}

	if err != nil {
		r.Log = Result{Status: ERR, Message: fmt.Sprintf("Ошибка чтения ответа на %s: %s", r.ReqId, err)}
		return
	}

	err = json.Unmarshal(body, r.Response)
	if err != nil {
		r.Log = Result{Status: ERR, Message: fmt.Sprintf("Ошибка десериализации %s: %s", r.ReqId, err)}
	}
}

func Sign(data, secret string, hash func() hash.Hash) string {
	mac := hmac.New(hash, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
