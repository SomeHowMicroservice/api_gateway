package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func PostRouter(rg *gin.RouterGroup, cfg *config.AppConfig, userClient userpb.UserServiceClient, postHandler *handler.PostHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	topic := rg.Group("/topics")
	{
		topic.GET("", postHandler.GetAllTopics)
	}

	admin := rg.Group("/admin", middleware.RequireAuth(accessName, secretKey, userClient), middleware.RequireMultiRoles([]string{common.RoleContributor}))
	{
		admin.POST("/topics", postHandler.CreateTopic)
		admin.GET("/topics", postHandler.GetAllTopicsAdmin)
		admin.GET("/topics/deleted", postHandler.GetDeletedTopics)
		admin.PUT("/topics/:id", postHandler.UpdateTopic)
		admin.DELETE("/topics/:id", postHandler.DeleteTopic)
		admin.DELETE("/topics", postHandler.DeleteTopics)
		admin.PUT("/topics/:id/restore", postHandler.RestoreTopic)
		admin.PUT("/topics/restore", postHandler.RestoreTopics)
		admin.DELETE("/topics/:id/permanent", postHandler.PermanentlyDeleteTopic)
		admin.DELETE("/topics/permanent", postHandler.PermanentlyDeleteTopics)
		admin.POST("/posts", postHandler.CreatePost)
		admin.GET("/posts", postHandler.GetAllPostsAdmin)
		admin.GET("/posts/:id", postHandler.GetPostByID)
		admin.PATCH("/posts/:id", postHandler.UpdatePost)
		admin.DELETE("/posts/:id", postHandler.DeletePost)
	}
}
