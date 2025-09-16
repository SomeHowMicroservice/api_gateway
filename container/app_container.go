package container

import (
	"github.com/SomeHowMicroservice/gateway/config"
	"github.com/SomeHowMicroservice/gateway/event"
	"github.com/SomeHowMicroservice/gateway/handler"
	"github.com/SomeHowMicroservice/gateway/initialization"
	"github.com/SomeHowMicroservice/gateway/socket"
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

func NewContainer(cs *initialization.GRPCClients, cfg *config.Config, hub *socket.Hub, manager *event.Manager) *Container {
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
