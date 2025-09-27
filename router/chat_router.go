package router

import (
	"github.com/SomeHowMicroservice/gateway/config"
	"github.com/SomeHowMicroservice/gateway/handler"
	"github.com/SomeHowMicroservice/gateway/security"
	userpb "github.com/SomeHowMicroservice/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func ChatRouter(rg *gin.RouterGroup, cfg *config.Config, userClient userpb.UserServiceClient, chatHandler *handler.ChatHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	chat := rg.Group("", security.RequireAuth(accessName, secretKey, userClient))
	{
		chat.GET("/me/conversations", security.RequireSingleRole(), chatHandler.MyConversation)
	}
}
