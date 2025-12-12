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
		auth.POST("/upload/tinjauan/:id", controllers.UploadExcelShopeeTinjauan)
		auth.POST("/upload/iklan/:id", controllers.UploadCsvShopeeIklan)

		auth.POST("/dashboard-tinjauan/total-penjualan", controllers.GetDashboardTinjauanTotalPenjualan)
		auth.POST("/dashboard-tinjauan/total-pesanan", controllers.GetDashboardTinjauanTotalPesanan)
		auth.POST("/dashboard-tinjauan/total-pengunjung", controllers.GetDashboardTinjauanTotalPengunjung)
		auth.POST("/dashboard-tinjauan/tren-penjualan", controllers.GetDashboardTinjauanTrenPenjualan)
		auth.POST("/dashboard-tinjauan/total-pesanan/in-week", controllers.GetDashboardTinjauanTotalPesananInWeek)
		auth.POST("/dashboard-tinjauan/convertion-rate", controllers.GetDashboardTinjauanConvertionRate)
		auth.POST("/dashboard-tinjauan/basket-size", controllers.GetDashboardTinjauanBasketSize)

		auth.POST("/dashboard-iklan/penjualan-iklan", controllers.GetDashboardIklanPenjualanIklan)
		auth.POST("/dashboard-iklan/biaya-iklan", controllers.GetDashboardIklanBiayaIklan)
		auth.POST("/dashboard-iklan/roas", controllers.GetDashboardIklanROAS)
		auth.POST("/dashboard-iklan/convertion-rate", controllers.GetDashboardIklanConvertionRateIklan)
		auth.POST("/dashboard-iklan/presentase-klik", controllers.GetDashboardIklanPresentaseKlik)
		auth.POST("/dashboard-iklan/dilihat", controllers.GetDashboardIklanDilihat)
		auth.POST("/dashboard-iklan/penjualan-dan-biaya", controllers.GetDashboardTotalPenjulandanBiaya)
		auth.POST("/dashboard-iklan/total-roas", controllers.GetDashboardTotalROAS)
		auth.POST("/dashboard-iklan/top-product", controllers.GetDashboardTopProduct)

		auth.GET("/history-data-upload", controllers.GetHistoryDataUpload)

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
