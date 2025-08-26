package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func ProductRouter(rg *gin.RouterGroup, cfg *config.AppConfig, userClient userpb.UserServiceClient, productHandler *handler.ProductHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	product := rg.Group("/products")
	{
		product.GET("/:slug", productHandler.GetProductBySlug)
	}

	category := rg.Group("/categories")
	{
		category.GET("/tree", productHandler.GetCategoryTree)
		category.GET("/:slug/products", productHandler.GetProductsByCategory)
		category.GET("/no-child", productHandler.GetCategoriesNoChild)
		category.GET("/no-product", productHandler.GetCategoriesNoProduct)
	}

	color := rg.Group("/colors")
	{
		color.GET("", productHandler.GetAllColors)
	}

	size := rg.Group("/sizes")
	{
		size.GET("", productHandler.GetAllSizes)
	}

	tag := rg.Group("/tags")
	{
		tag.GET("", productHandler.GetAllTags)
	}

	admin := rg.Group("/admin", middleware.RequireAuth(accessName, secretKey, userClient), middleware.RequireMultiRoles([]string{common.RoleAdmin}))
	{
		admin.POST("/categories", productHandler.CreateCategory)
		admin.GET("/categories", productHandler.GetAllCategoriesAdmin)
		admin.GET("/categories/:id", productHandler.CategoryAdminDetails)
		admin.PUT("/categories/:id", productHandler.UpdateCategory)
		admin.DELETE("/categories/permanent", productHandler.PermanentlyDeleteCategories)
		admin.DELETE("/categories/:id/permanent", productHandler.PermanentlyDeleteCategory)
		admin.POST("/products", productHandler.CreateProduct)
		admin.GET("/products", productHandler.GetAllProductsAdmin)
		admin.DELETE("/products/:id", productHandler.DeleteProduct)
		admin.DELETE("/products", productHandler.DeleteProducts)
		admin.DELETE("/products/permanent", productHandler.PermanentlyDeleteProducts)
		admin.DELETE("/products/:id/permanent", productHandler.PermanentlyDeleteProduct)
		admin.GET("/products/deleted", productHandler.GetDeletedProducts)
		admin.GET("/products/:id", productHandler.GetProductByID)
		admin.GET("/products/:id/deleted", productHandler.GetDeletedProductByID)
		admin.PUT("/products/restore", productHandler.RestoreProducts)
		admin.PATCH("/products/:id", productHandler.UpdateProduct)
		admin.PATCH("/products/:id/restore", productHandler.RestoreProduct)
		admin.POST("/colors", productHandler.CreateColor)
		admin.GET("/colors", productHandler.GetAllColorsAdmin)
		admin.GET("/colors/deleted", productHandler.GetDeletedColors)
		admin.PUT("/colors/:id", productHandler.UpdateColor)
		admin.PUT("/colors/restore", productHandler.RestoreColors)
		admin.PATCH("/colors/:id/restore", productHandler.RestoreColor)
		admin.DELETE("/colors", productHandler.DeleteColors)
		admin.DELETE("/colors/:id", productHandler.DeleteColor)
		admin.DELETE("/colors/permanent", productHandler.PermanentlyDeleteColors)
		admin.DELETE("/colors/:id/permanent", productHandler.PermanentlyDeleteColor)
		admin.POST("/sizes", productHandler.CreateSize)
		admin.GET("/sizes", productHandler.GetAllSizesAdmin)
		admin.GET("/sizes/deleted", productHandler.GetDeletedSizes)
		admin.PUT("/sizes/restore", productHandler.RestoreSizes)
		admin.PATCH("/sizes/:id/restore", productHandler.RestoreSize)
		admin.PUT("/sizes/:id", productHandler.UpdateSize)
		admin.DELETE("/sizes", productHandler.DeleteSizes)
		admin.DELETE("/sizes/:id", productHandler.DeleteSize)
		admin.DELETE("/sizes/permanent", productHandler.PermanentlyDeleteSizes)
		admin.DELETE("/sizes/:id/permanent", productHandler.PermanentlyDeleteSize)
		admin.POST("/tags", productHandler.CreateTag)
		admin.GET("/tags", productHandler.GetAllTagsAdmin)
		admin.GET("/tags/deleted", productHandler.GetDeletedTags)
		admin.PUT("/tags/restore", productHandler.RestoreTags)
		admin.PATCH("/tags/:id/restore", productHandler.RestoreTag)
		admin.PUT("/tags/:id", productHandler.UpdateTag)
		admin.DELETE("/tags", productHandler.DeleteTags)
		admin.DELETE("/tags/:id", productHandler.DeleteTag)
		admin.DELETE("/tags/permanent", productHandler.PermanentlyDeleteTags)
		admin.DELETE("/tags/:id/permanent", productHandler.PermanentlyDeleteTag)
	}
}
