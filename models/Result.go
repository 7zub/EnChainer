package models

type Result struct {
	Status  Status
	Message string
	Any     any `gorm:"-"`
}

type Status string

const (
	OK   Status = "OK"
	INFO Status = "INFO"
	WAR  Status = "WAR"
	ERR  Status = "ERR"
)

func (r Result) GetStatus() any {
	return r.Status
}

func (r Result) GetMessage() string {
	return r.Message
}
