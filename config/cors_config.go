package config

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	allowOrigins  = []string{"http://localhost:3000", "http://127.0.0.1:5500"}
	allowMethods  = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	allowHeaders  = []string{"Origin", "Content-Type", "Content-Length", "Authorization", "Accept", "Accept-Language"}
	exposeHeaders = []string{"Content-Length", "X-Request-Id"}
)

func CORSConfig(r *gin.Engine) {
	config := cors.Config{
		AllowOrigins:        allowOrigins,
		AllowMethods:        allowMethods,
		AllowHeaders:        allowHeaders,
		ExposeHeaders:       exposeHeaders,
		AllowCredentials:    true,
		AllowWebSockets:     true,
		AllowFiles:          true,
		AllowPrivateNetwork: true,
		MaxAge:              12 * time.Hour,
	}
	r.Use(cors.New(config))
}
