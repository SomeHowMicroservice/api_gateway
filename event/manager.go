package event

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
)

type Manager struct {
	Users      map[string]*User
	Register   chan *User
	Unregister chan *User
	Broadcast  chan []byte
	Mutex      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		Users:      make(map[string]*User),
		Register:   make(chan *User),
		Unregister: make(chan *User),
		Broadcast:  make(chan []byte),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case user := <-m.Register:
			m.Mutex.Lock()
			m.Users[user.ID] = user
			m.Mutex.Unlock()
			log.Printf("Người dùng %s đã kết nối", user.ID)

		case user := <-m.Unregister:
			m.Mutex.Lock()
			if _, ok := m.Users[user.ID]; ok {
				delete(m.Users, user.ID)
				close(user.Send)
				close(user.Done)
			}
			m.Mutex.Unlock()
			log.Printf("Người dùng %s đã ngắt kết nối", user.ID)

		case message := <-m.Broadcast:
			m.Mutex.RLock()
			for _, user := range m.Users {
				select {
				case user.Send <- message:
				default:
					delete(m.Users, user.ID)
					close(user.Send)
					close(user.Done)
				}
			}
			m.Mutex.RUnlock()
		}
	}
}

func (m *Manager) SendToUser(userID string, event *common.SSEEvent) {
	data, _ := json.Marshal(event)

	m.Mutex.RLock()
	for _, user := range m.Users {
		if user.UserID == userID {
			select {
			case user.Send <- data:
			default:
				m.Mutex.RUnlock()
				m.Mutex.Lock()
				if _, ok := m.Users[user.ID]; ok {
					delete(m.Users, user.ID)
					close(user.Send)
					close(user.Done)
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
