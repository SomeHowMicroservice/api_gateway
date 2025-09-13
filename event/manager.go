package event

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
)

type Manager struct {
	Clients      map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	Mutex      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		Clients:      make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.Register:
			m.Mutex.Lock()
			m.Clients[client.ID] = client
			m.Mutex.Unlock()
			log.Printf("Người dùng %s đã kết nối", client.ID)

		case client := <-m.Unregister:
			m.Mutex.Lock()
			if _, ok := m.Clients[client.ID]; ok {
				delete(m.Clients, client.ID)
				close(client.Send)
				close(client.Done)
			}
			m.Mutex.Unlock()
			log.Printf("Người dùng %s đã ngắt kết nối", client.ID)

		case message := <-m.Broadcast:
			m.Mutex.RLock()
			for _, client := range m.Clients {
				select {
				case client.Send <- message:
				default:
					delete(m.Clients, client.ID)
					close(client.Send)
					close(client.Done)
				}
			}
			m.Mutex.RUnlock()
		}
	}
}

func (m *Manager) SendToUser(userID string, event *common.SSEEvent) {
	data, _ := json.Marshal(event)

	m.Mutex.RLock()
	for _, client := range m.Clients {
		if client.UserID == userID {
			select {
			case client.Send <- data:
			default:
				m.Mutex.RUnlock()
				m.Mutex.Lock()
				if _, ok := m.Clients[client.ID]; ok {
					delete(m.Clients, client.ID)
					close(client.Send)
					close(client.Done)
				}
				m.Mutex.Unlock()
				m.Mutex.RLock()
			}
		}
	}
	m.Mutex.RUnlock()
}

func (m *Manager) BroadcastToAll(event *common.SSEEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling SSE event: %v", err)
		return
	}

	select {
	case m.Broadcast <- data:
	default:
		log.Println("Broadcast channel is full")
	}
}
