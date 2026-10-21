package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kan/mock-api/internal/config"
	"github.com/kan/mock-api/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ ไม่พบไฟล์ .env ใช้ Environment Variables จากระบบ")
	}

	config.ConnectDatabase()

	r := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400,
	}
	r.Use(cors.New(corsConfig))

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
