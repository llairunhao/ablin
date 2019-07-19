package dao

type GGStatusCode int

const (
	StatusOk GGStatusCode = 10000 //成功

	StatusParameterError = 20001
	StatusMySqlError     = 20002
)
