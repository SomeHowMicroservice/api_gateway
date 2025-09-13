package event

import "github.com/google/uuid"

type Client struct {
	ID     string
	UserID string
	Send   chan []byte
	Done   chan bool
}

func NewClient(userID string) *Client {
	return &Client{
		ID: uuid.NewString(),
		UserID: userID,
		Send: make(chan []byte, 256),
		Done: make(chan bool),
	}
}