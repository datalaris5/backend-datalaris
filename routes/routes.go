package routes

import (
	"go-datalaris/controllers"
	"go-datalaris/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/umkm/v1/api")

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

		auth.POST("/marketplaces", controllers.CreateStore)
		auth.PUT("/marketplaces", controllers.UpdateStore)
		auth.GET("/marketplaces/:id", controllers.GetStoreById)
	}
}
