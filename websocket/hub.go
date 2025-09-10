package websocket

import "log"

type Hub struct {
	Conversations map[string]map[*Client]bool
	Broadcast     chan *Message
	Register      chan *Client
	Unregister    chan *Client
}

type Message struct {
	ConversationID string
	Content        []byte
}

func NewHub() *Hub {
	return &Hub{
		Conversations: make(map[string]map[*Client]bool),
		Broadcast:     make(chan *Message),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if h.Conversations[client.ConversationID] == nil {
				h.Conversations[client.ConversationID] = make(map[*Client]bool)
			}

			h.Conversations[client.ConversationID][client] = true
			log.Printf("Người dùng đã tham gia cuộc trò chuyện %s. Tổng số người tham gia trong cuộc trò chuyện: %d",
				client.ConversationID, len(h.Conversations[client.ConversationID]))

		case client := <-h.Unregister:
			if clients, ok := h.Conversations[client.ConversationID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)

					if len(clients) == 0 {
						delete(h.Conversations, client.ConversationID)
						log.Printf("Cuộc trò chuyện %s đã được xóa (không có người tham gia)", client.ConversationID)
					} else {
						log.Printf("Người dùng đã rời khỏi cuộc trò chuyện %s. Số người còn lại: %d",
							client.ConversationID, len(clients))
					}
				}
			}

		case message := <-h.Broadcast:
			if clients, ok := h.Conversations[message.ConversationID]; ok {
				for client := range clients {
					select {
					case client.Send <- message.Content:
					default:
						close(client.Send)
						delete(clients, client)

						if len(clients) == 0 {
							delete(h.Conversations, message.ConversationID)
						}
					}
				}
			}
		}
	}
}
