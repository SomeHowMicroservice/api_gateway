package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func SSERouter(rg *gin.RouterGroup, cfg *config.Config, userClient userpb.UserServiceClient, sseHandler *handler.SSEHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	rg.GET("/sse", middleware.OptionalAuth(accessName, secretKey, userClient), sseHandler.HandleSSE)
}
