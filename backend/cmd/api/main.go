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
	// Load .env if present (on cloud hosts, env vars are injected directly)
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ ไม่พบไฟล์ .env ใช้ Environment Variables จากระบบ")
	}

	//connect db
	config.ConnectDatabase()

	// gin
	r := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false, // Must be false when AllowAllOrigins is true
		MaxAge:           86400,
	}
	r.Use(cors.New(corsConfig))

	//endpoint
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
