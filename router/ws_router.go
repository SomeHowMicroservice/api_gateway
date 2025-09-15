package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func WSRouter(rg *gin.RouterGroup, cfg *config.Config, userClient userpb.UserServiceClient, wsHandler *handler.WSHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	rg.GET("/ws", middleware.RequireAuth(accessName, secretKey, userClient), wsHandler.HandleWS)
}
