package router

import (
	"github.com/SomeHowMicroservice/gateway/config"
	"github.com/SomeHowMicroservice/gateway/handler"
	"github.com/SomeHowMicroservice/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func ChatRouter(rg *gin.RouterGroup, cfg *config.Config, userClient userpb.UserServiceClient, chatHandler *handler.ChatHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	chat := rg.Group("", middleware.RequireAuth(accessName, secretKey, userClient))
	{
		chat.GET("/me/conversations", middleware.RequireSingleRole(), chatHandler.MyConversation)
	}
}
