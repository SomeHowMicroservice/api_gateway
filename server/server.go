package server

import (
	"context"
	"fmt"
	"log"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/container"
	"github.com/SomeHowMicroservice/shm-be/gateway/initialization"
	"github.com/SomeHowMicroservice/shm-be/gateway/websocket"
)

var (
	authAddr    = "localhost:8081"
	userAddr    = "localhost:8082"
	productAddr = "localhost:8083"
	postAddr    = "localhost:8084"
	chatAddr    = "localhost:8085"
)

type Server struct {
	cfg          *config.AppConfig
	clients      *initialization.GRPCClients
	appContainer *container.Container
	hub          *websocket.Hub
}

func NewServer(cfg *config.AppConfig) (*Server, error) {
	authAddr = fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.Services.AuthPort)
	userAddr = fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.Services.UserPort)
	productAddr = fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.Services.ProductPort)
	postAddr = fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.Services.PostPort)
	chatAddr = fmt.Sprintf("%s:%d", cfg.App.ServerHost, cfg.Services.ChatPort)

	ca := &common.ClientAddresses{
		AuthAddr:    authAddr,
		UserAddr:    userAddr,
		ProductAddr: productAddr,
		PostAddr:    postAddr,
		ChatAddr:    chatAddr,
	}

	clients, err := initialization.InitClients(ca)
	if err != nil {
		return nil, fmt.Errorf("init clients thất bại: %w", err)
	}

	hub := websocket.NewHub()
	go hub.Run()

	appContainer := container.NewContainer(clients, cfg, hub)

	return &Server{
		cfg:          cfg,
		clients:      clients,
		appContainer: appContainer,
		hub:          hub,
	}, nil
}

func (s *Server) Start() error {
	return RunHTTPServer(s.cfg, s.clients, s.appContainer)
}

func (s *Server) Shutdown(ctx context.Context) {
	log.Println("Đang shutdown service...")

	if s.clients != nil {
		s.clients.Close()
	}

	if s.hub != nil {
		for _, clients := range s.hub.Conversations {
			for client := range clients {
				close(client.Send)
			}
		}
	}

	log.Println("Shutdown service thành công")
}
