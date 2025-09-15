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
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewHttpServer(cfg *config.Config, clients *initialization.GRPCClients, hub *socket.Hub, manager *event.Manager) (*http.Server, error) {
	appContainer := container.NewContainer(clients, cfg, hub, manager)

	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil, fmt.Errorf("thiết lập Proxy thất bại: %w", err)
	}

	corsConfig := cors.Config{
		AllowOrigins:        cfg.App.CORS.AllowOrigins,
		AllowMethods:        cfg.App.CORS.AllowMethods,
		AllowHeaders:        cfg.App.CORS.AllowHeaders,
		ExposeHeaders:       cfg.App.CORS.ExposeHeaders,
		AllowCredentials:    cfg.App.CORS.AllowCredentials,
		AllowWebSockets:     cfg.App.CORS.AllowWebSockets,
		AllowFiles:          cfg.App.CORS.AllowFiles,
		AllowPrivateNetwork: cfg.App.CORS.AllowPrivateNetwork,
		MaxAge:              cfg.App.CORS.MaxAge * time.Hour,
	}

	r.Use(cors.New(corsConfig))

	api := r.Group("/api/v1")
	router.AuthRouter(api, cfg, clients.UserClient, appContainer.Auth.Handler)
	router.UserRouter(api, cfg, clients.UserClient, appContainer.User.Handler)
	router.ProductRouter(api, cfg, clients.UserClient, appContainer.Product.Handler)
	router.PostRouter(api, cfg, clients.UserClient, appContainer.Post.Handler)
	router.ChatRouter(api, cfg, clients.UserClient, appContainer.Chat.Handler)
	router.SSERouter(api, cfg, clients.UserClient, appContainer.SSEHandler)
	router.WSRouter(api, cfg, clients.UserClient, appContainer.WSHandler)

	addr := fmt.Sprintf(":%d", cfg.App.HttpPort)

	httpServer := &http.Server{
		Addr:           addr,
		Handler:        r,
		IdleTimeout:    cfg.App.Http.IdleTimeout * time.Second,
		ReadTimeout:    cfg.App.Http.ReadTimeout * time.Second,
		WriteTimeout:   cfg.App.Http.ReadTimeout * time.Second,
		MaxHeaderBytes: cfg.App.Http.MaxHeaderBytes * 1024 * 1024,
	}

	return httpServer, nil
}
