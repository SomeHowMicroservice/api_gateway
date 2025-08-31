package main

import (
	"fmt"
	"log"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/container"
	"github.com/SomeHowMicroservice/shm-be/gateway/initialization"
	"github.com/SomeHowMicroservice/shm-be/gateway/server"
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

	server.RunHTTPServer(cfg, clients, appContainer)
}
