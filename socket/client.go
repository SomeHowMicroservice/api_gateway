package socket

import (
	"context"
	"encoding/json"
	"log"
	"time"

	chatpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/chat"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/status"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

type Client struct {
	Hub            *Hub
	Conn           *websocket.Conn
	Send           chan []byte
	UserID         string
	UserRole       string
	ConversationID string
}

func NewClient(hub *Hub, conn *websocket.Conn, userID, userRole, conversationID string) *Client {
	return &Client{
		hub,
		conn,
		make(chan []byte, 256),
		userID,
		userRole,
		conversationID,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		messageType, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var file []byte
		var content string
		switch messageType {
		case websocket.BinaryMessage:
			file = data
		case websocket.TextMessage:
			content = string(data)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := c.Hub.ChatClient.CreateMessage(ctx, &chatpb.CreateMessageRequest{
			ConversationId: c.ConversationID,
			SenderId:       c.UserID,
			SenderRole:     c.UserRole,
			Content:        &content,
			FileData:       file,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				log.Println("gRPC error ", st.Message())
				break
			}
			log.Println("Chat service err ", err)
			break
		}

		messageBytes, _ := json.Marshal(res)

		message := &Message{
			ConversationID: c.ConversationID,
			Content:        messageBytes,
		}

		c.Hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
