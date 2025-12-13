package main

import (
	"fmt"
	"go-datalaris/config"
	"go-datalaris/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	db := config.DB

	r := gin.Default()

	// r.LoadHTMLGlob("templates/*.html")

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://staging-cms.thesanur.id",
			"https://staging.thesanur.id",
			"https://cms.thesanur.id",
			"https://thesanur.id",
			"http://localhost:3000",
			"http://localhost:5173",
			"https://staging-cms.injourneyhospitality.id",
			"https://hotel-management-cms.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(r, db)

	// seed.Seed(db)

	r.Run(fmt.Sprintf(":%s", config.AppPort))
}
