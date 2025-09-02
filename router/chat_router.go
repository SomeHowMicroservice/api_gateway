package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func ChatRouter(rg *gin.RouterGroup, cfg *config.AppConfig, userClient userpb.UserServiceClient, chatHandler *handler.ChatHandler) {
	chat := rg.Group("/chats")
	{
		chat.GET("", chatHandler.TestConnect)
		chat.GET("/ws", chatHandler.ServeWs)
	}
}
