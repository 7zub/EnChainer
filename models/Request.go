package models

import "awesomeProject/controllers"

type IParams interface {
	GetParams(ccy string) Request
}

type Request struct {
	Exchange int
	Url      string
	Params   interface{}
	Response interface{}
}

func (r Request) SendRequest() {
	controllers.ApiGetBook1(r)
}
