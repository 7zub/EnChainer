package models

type Result struct {
	Status  Status
	Message string
	Any     any
}

type Status string

const (
	OK   Status = "OK"
	INFO Status = "INFO"
	WAR  Status = "WAR"
	ERR  Status = "ERR"
)
