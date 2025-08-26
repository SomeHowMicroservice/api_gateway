package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/initialization"
)

type Container struct {
	Auth    *AuthContainer
	User    *UserContainer
	Product *ProductContainer
	Post    *PostContainer
	Chat    *ChatContainer
}

func NewContainer(cs *initialization.GRPCClients, cfg *config.AppConfig) *Container {
	auth := NewAuthContainer(cs.AuthClient, cfg)
	user := NewUserContainer(cs.UserClient)
	product := NewProductHandler(cs.ProductClient)
	post := NewPostContainer(cs.PostClient)
	chat := NewChatContainer(cs.ChatClient)
	return &Container{
		auth,
		user,
		product,
		post,
		chat,
	}
}
