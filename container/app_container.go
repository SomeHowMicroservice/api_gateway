package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/initialization"
	"github.com/SomeHowMicroservice/shm-be/gateway/websocket"
)

type Container struct {
	Auth    *AuthContainer
	User    *UserContainer
	Product *ProductContainer
	Post    *PostContainer
	Chat    *ChatContainer
}

func NewContainer(cs *initialization.GRPCClients, cfg *config.AppConfig, hub *websocket.Hub) *Container {
	auth := NewAuthContainer(cs.AuthClient, cfg)
	user := NewUserContainer(cs.UserClient)
	product := NewProductHandler(cs.ProductClient)
	post := NewPostContainer(cs.PostClient)
	chat := NewChatContainer(cs.ChatClient, hub)
	return &Container{
		auth,
		user,
		product,
		post,
		chat,
	}
}
