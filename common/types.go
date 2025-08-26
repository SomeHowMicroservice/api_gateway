package common

type ApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ClientAddresses struct {
	UserAddr    string
	AuthAddr    string
	ProductAddr string
	PostAddr    string
	ChatAddr    string
}
