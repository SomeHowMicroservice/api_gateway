package common

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ClientAddresses struct {
	UserAddr    string
	AuthAddr    string
	ProductAddr string
	PostAddr    string
	ChatAddr    string
}

type ImageUploadedEvent struct {
	Service string `json:"service"`
	UserID  string `json:"user_id"`
}

type SSEEvent struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}
