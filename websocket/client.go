package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Thời gian chờ để ghi tin nhắn đến peer.
	writeWait = 10 * time.Second

	// Thời gian chờ để đọc tin nhắn tiếp theo từ peer.
	pongWait = 60 * time.Second

	// Gửi pings đến peer với khoảng thời gian này. Phải nhỏ hơn pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Kích thước buffer tối đa cho tin nhắn.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	ConversationID string
}

func NewClient(hub *Hub, conn *websocket.Conn, conversationID string) *Client {
	return &Client{
		Hub:            hub,
		Conn:           conn,
		Send:           make(chan []byte, 256),
		ConversationID: conversationID,
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
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		
		// Tạo Message với ConversationID
		message := &Message{
			ConversationID: c.ConversationID,
			Content:        messageBytes,
		}
		
		// Gửi message đến Hub để broadcast
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
				// Hub đã đóng channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Thêm các tin nhắn đang chờ trong queue vào tin nhắn hiện tại.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
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
