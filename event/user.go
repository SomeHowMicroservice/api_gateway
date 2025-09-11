package event

import "github.com/google/uuid"

type User struct {
	ID     string
	UserID string
	Send   chan []byte
	Done   chan bool
}

func NewUserReceiverEvent(userID string) *User {
	return &User{
		ID: uuid.NewString(),
		UserID: userID,
		Send: make(chan []byte, 256),
		Done: make(chan bool),
	}
}