package container

import (
	"github.com/SomeHowMicroservice/gateway/config"
	"github.com/SomeHowMicroservice/gateway/handler"
	authpb "github.com/SomeHowMicroservice/gateway/protobuf/auth"
)

type AuthContainer struct {
	Handler *handler.AuthHandler
}

func NewAuthContainer(authClient authpb.AuthServiceClient, cfg *config.Config) *AuthContainer {
	handler := handler.NewAuthHandler(authClient, cfg)
	return &AuthContainer{handler}
}
