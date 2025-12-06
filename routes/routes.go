package routes

import (
	"go-datalaris/controllers"
	"go-datalaris/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/datalaris/v1/api")

	// --- PUBLIC ---
	api.POST("/login", controllers.Login)

	// --- PROTECTED ---
	auth := api.Group("/admin")
	auth.Use(middlewares.AuthMiddlewareWithDB())
	{
		auth.POST("/upload", controllers.UploadExcel)

		auth.POST("/store", controllers.CreateStore)
		auth.PUT("/store", controllers.UpdateStore)
		auth.GET("/store/:id", controllers.GetStoreById)
		auth.POST("/store/page", controllers.GetStoreByPage)
		auth.DELETE("/store/:id", controllers.SoftDeleteStore)
		auth.POST("/store/:id/status/:status", controllers.ActiveInactiveStore)

		auth.POST("/marketplaces", controllers.CreateMarketplace)
		auth.PUT("/marketplaces", controllers.UpdateMarketplace)
		auth.GET("/marketplaces/:id", controllers.GetMarketplaceById)
		auth.POST("/marketplaces/page", controllers.GetMarketplaceByPage)
		auth.DELETE("/marketplaces/:id", controllers.SoftDeleteMarketplace)
		auth.POST("/marketplaces/:id/status/:status", controllers.ActiveInactiveMarketplace)

		//LOV
		auth.GET("marketplaces/lov", controllers.GetLovMarketplace)
		auth.GET("store/lov", controllers.GetLovStore)

	}
}
