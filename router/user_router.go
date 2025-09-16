package router

import (
	"github.com/SomeHowMicroservice/gateway/config"
	"github.com/SomeHowMicroservice/gateway/handler"
	"github.com/SomeHowMicroservice/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func UserRouter(rg *gin.RouterGroup, cfg *config.Config, userClient userpb.UserServiceClient, userHandler *handler.UserHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	user := rg.Group("", middleware.RequireAuth(accessName, secretKey, userClient))
	{
		user.PATCH("/profiles/:id", userHandler.UpdateProfile)
		user.GET("/me/measurements", userHandler.MyMeasurements)
		user.PATCH("/measurements/:id", userHandler.UpdateMeasurement)
		user.GET("/me/addresses", userHandler.MyAddresses)
		user.POST("/me/addresses", userHandler.CreateMyAddress)
		user.PUT("/addresses/:id", userHandler.UpdateAddress)
		user.DELETE("/addresses/:id", userHandler.DeleteAddress)
	}
}
