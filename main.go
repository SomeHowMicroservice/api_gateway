package main

import (
	"fmt"
	"log"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/container"
	"github.com/SomeHowMicroservice/shm-be/gateway/initialization"
	"github.com/SomeHowMicroservice/shm-be/gateway/router"
	"github.com/gin-gonic/gin"
)

var (
	authAddr    = "localhost:8081"
	userAddr    = "localhost:8082"
	productAddr = "localhost:8083"
	postAddr    = "localhost:8084"
	chatAddr    = "localhost:8085"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình Gateway thất bại: %v", err)
	}

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
		log.Fatalf("Kết nối tới các dịch vụ khác thất bại: %v", err)
	}
	defer clients.Close()

	appContainer := container.NewContainer(clients, cfg)

	r := gin.Default()
	if err = r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("Thiết lập Proxy thất bại: %v", err)
	}

	config.CORSConfig(r)

	api := r.Group("/api/v1")
	router.AuthRouter(api, cfg, clients.UserClient, appContainer.Auth.Handler)
	router.UserRouter(api, cfg, clients.UserClient, appContainer.User.Handler)
	router.ProductRouter(api, cfg, clients.UserClient, appContainer.Product.Handler)
	router.PostRouter(api, cfg, clients.UserClient, appContainer.Post.Handler)
	router.ChatRouter(api, cfg, clients.UserClient, appContainer.Chat.Handler)

	r.Run(fmt.Sprintf(":%d", cfg.App.GRPCPort))
}
