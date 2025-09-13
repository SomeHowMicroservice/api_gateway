package socket

import (
	"encoding/json"
	"log"
	"time"

	chatpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/chat"
)

type Hub struct {
	Conversations map[string]map[*Client]bool
	UserConnCount map[string]int
	UserOnline    map[string]bool
	UserLastSeen  map[string]time.Time
	Broadcast     chan *Message
	Register      chan *Client
	Unregister    chan *Client
	ChatClient    chatpb.ChatServiceClient
}

type Message struct {
	ConversationID string
	Content        []byte
}

func NewHub(chatClient chatpb.ChatServiceClient) *Hub {
	return &Hub{
		make(map[string]map[*Client]bool),
		make(map[string]int),
		make(map[string]bool),
		make(map[string]time.Time),
		make(chan *Message),
		make(chan *Client),
		make(chan *Client),
		chatClient,
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

			h.UserConnCount[client.UserID]++
			if h.UserConnCount[client.UserID] == 1 {
				h.UserOnline[client.UserID] = true
				h.notifyUserStatus(client.ConversationID, client.UserID, true)
			}

			log.Printf("Tổng số kết nối trong cuộc trò chuyện %s là %d", client.ConversationID, len(h.Conversations[client.ConversationID]))

		case client := <-h.Unregister:
			if clients, ok := h.Conversations[client.ConversationID]; ok {
				if _, exists := clients[client]; exists {
					h.UserConnCount[client.UserID]--
					if h.UserConnCount[client.UserID] <= 0 {
						h.UserOnline[client.UserID] = false
						delete(h.UserConnCount, client.UserID)

						h.UserLastSeen[client.UserID] = time.Now()

						h.notifyUserStatus(client.ConversationID, client.UserID, false)
					}

					delete(clients, client)
					close(client.Send)

					if len(clients) == 0 {
						delete(h.Conversations, client.ConversationID)
						log.Printf("Cuộc trò chuyện %s đã được xóa (không có người tham gia)", client.ConversationID)
					} else {
						log.Printf("Có thiết bị ngắt kết nối, số kết nối còn lại trong cuộc trò chuyện %s là %d", client.ConversationID, len(clients))
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

func (h *Hub) notifyUserStatus(conversationID, userID string, online bool) {
	statusMessage := map[string]any{
		"type":    "user_status",
		"user_id": userID,
		"online":  online,
	}

	if !online {
		if t, ok := h.UserLastSeen[userID]; ok {
			statusMessage["last_seen"] = t.Format(time.RFC3339)
		}
	}

	messageBytes, err := json.Marshal(statusMessage)
	if err != nil {
		log.Println("Lỗi khi mã hóa  tin nhắn user_status", err.Error())
		return
	}

	if clients, ok := h.Conversations[conversationID]; ok {
		for client := range clients {
			select {
			case client.Send <- messageBytes:
			default:
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}
