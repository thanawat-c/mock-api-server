package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kan/mock-api/internal/config"
	"github.com/kan/mock-api/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ ไม่พบไฟล์ .env")
	}

	//connect db
	config.ConnectDatabase()

	// gin
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization", "Origin", "x-mock-status"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge: 86400,
	}))

	//endpoint
	routes.SetupRoutes(r)

	r.Run(":8080")
}