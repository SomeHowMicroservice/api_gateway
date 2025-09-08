package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func ChatRouter(rg *gin.RouterGroup, cfg *config.AppConfig, userClient userpb.UserServiceClient, chatHandler *handler.ChatHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	chat := rg.Group("", middleware.RequireAuth(accessName, secretKey, userClient))
	{
		chat.GET("/me/conversations", middleware.RequireSingleRole(), chatHandler.MyConversation)
		chat.GET("/ws", chatHandler.ServeWs)
	}
}
