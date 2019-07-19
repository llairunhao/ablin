package dao

type GGResponse struct {
	Code    GGStatusCode `json:"code,omitempty"`
	Message string       `json:"message,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
}
