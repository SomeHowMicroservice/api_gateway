package container

import (
	"github.com/SomeHowMicroservice/gateway/handler"
	userpb "github.com/SomeHowMicroservice/gateway/protobuf/user"
)

type UserContainer struct {
	Handler *handler.UserHandler
}

func NewUserContainer(userClient userpb.UserServiceClient) *UserContainer {
	handler := handler.NewUserHandler(userClient)
	return &UserContainer{handler}
}