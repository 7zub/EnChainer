package models

import (
	"crypto/hmac"
	"crypto/sha256"
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

type IParams interface {
	GetParams(t any) *Request
}

type Request struct {
	ReqId       string `gorm:"primaryKey"`
	ReqType     string
	SignType    func() hash.Hash `gorm:"-"`
	Url         string
	Head        http.Header `gorm:"-"`
	Params      IParams     `gorm:"-"`
	Response    IResponse   `gorm:"-"`
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
	reqId := fmt.Sprintf("%02d/%02d%02d%02d-%d",
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

	switch r.ReqType {
	case "Trade", "Balance":
		rq.Method = "POST"
		rq.Header = r.Head
		signature := sign(q.Encode(), Conf.SecretKey)
		q.Add("signature", signature)
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

func sign(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
