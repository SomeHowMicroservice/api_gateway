package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/gin-gonic/gin"
)

func SSERouter(rg *gin.RouterGroup, sseHandler *handler.SSEHandler) {
	rg.GET("/sse", sseHandler.HandleSSE)
}
