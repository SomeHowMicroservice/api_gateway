package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/event"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/initialization"
	"github.com/SomeHowMicroservice/shm-be/gateway/socket"
)

type Container struct {
	Auth       *AuthContainer
	User       *UserContainer
	Product    *ProductContainer
	Post       *PostContainer
	Chat       *ChatContainer
	SSEHandler *handler.SSEHandler
	WSHandler  *handler.WSHandler
}

func NewContainer(cs *initialization.GRPCClients, cfg *config.AppConfig, hub *socket.Hub, manager *event.Manager) *Container {
	auth := NewAuthContainer(cs.AuthClient, cfg)
	user := NewUserContainer(cs.UserClient)
	product := NewProductHandler(cs.ProductClient)
	post := NewPostContainer(cs.PostClient)
	chat := NewChatContainer(cs.ChatClient)
	sseHandler := handler.NewSSEHandler(manager, cs.UserClient)
	wsHandler := handler.NewWSHandler(hub)
	return &Container{
		auth,
		user,
		product,
		post,
		chat,
		sseHandler,
		wsHandler,
	}
}
