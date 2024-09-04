package models

type Request struct {
	Url      string
	Currency string
	Params   interface{}
}
