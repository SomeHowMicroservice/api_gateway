package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/container"
	"github.com/SomeHowMicroservice/shm-be/gateway/event"
	"github.com/SomeHowMicroservice/shm-be/gateway/initialization"
	"github.com/SomeHowMicroservice/shm-be/gateway/router"
	"github.com/SomeHowMicroservice/shm-be/gateway/socket"
	"github.com/gin-gonic/gin"
)

func NewHttpServer(cfg *config.AppConfig, clients *initialization.GRPCClients, hub *socket.Hub, manager *event.Manager) (*http.Server, error) {
	appContainer := container.NewContainer(clients, cfg, hub, manager)

	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil, fmt.Errorf("thiết lập Proxy thất bại: %w", err)
	}

	config.CORSConfig(r)

	api := r.Group("/api/v1")
	router.AuthRouter(api, cfg, clients.UserClient, appContainer.Auth.Handler)
	router.UserRouter(api, cfg, clients.UserClient, appContainer.User.Handler)
	router.ProductRouter(api, cfg, clients.UserClient, appContainer.Product.Handler)
	router.PostRouter(api, cfg, clients.UserClient, appContainer.Post.Handler)
	router.ChatRouter(api, cfg, clients.UserClient, appContainer.Chat.Handler)
	router.SSERouter(api, appContainer.SSEHandler)
	router.WSRouter(api, cfg, clients.UserClient, appContainer.WSHandler)

	addr := fmt.Sprintf(":%d", cfg.App.HttpPort)

	httpServer := &http.Server{
		Addr:           addr,
		Handler:        r,
		IdleTimeout:    time.Minute,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   0,
		MaxHeaderBytes: 1 << 20,
	}

	return httpServer, nil
}
