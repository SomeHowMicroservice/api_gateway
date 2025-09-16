package router

import (
	"github.com/SomeHowMicroservice/gateway/config"
	"github.com/SomeHowMicroservice/gateway/handler"
	"github.com/SomeHowMicroservice/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func SSERouter(rg *gin.RouterGroup, cfg *config.Config, userClient userpb.UserServiceClient, sseHandler *handler.SSEHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	rg.GET("/sse", middleware.OptionalAuth(accessName, secretKey, userClient), sseHandler.HandleSSE)
}
