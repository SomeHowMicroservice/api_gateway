package router

import (
	"github.com/SomeHowMicroservice/gateway/config"
	"github.com/SomeHowMicroservice/gateway/handler"
	"github.com/SomeHowMicroservice/gateway/security"
	userpb "github.com/SomeHowMicroservice/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func WSRouter(rg *gin.RouterGroup, cfg *config.Config, userClient userpb.UserServiceClient, wsHandler *handler.WSHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	rg.GET("/ws", security.RequireAuth(accessName, secretKey, userClient), wsHandler.HandleWS)
}
