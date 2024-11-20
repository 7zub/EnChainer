package models

type Result struct {
	Status string
	//Source  string
	Message string
}

const (
	INFO = "INFO"
	WAR  = "WAR"
	ERR  = "ERR"
)
