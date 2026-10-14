package config

import (
	"log"
	"os"

	"github.com/kan/mock-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("❌ ไม่พบค่า DB_URL ในไฟล์ .env")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ เชื่อมต่อ Database ไม่สำเร็จ: ", err)
	}

	err = database.AutoMigrate(&models.Endpoint{}, &models.MockRecord{})
	if	err != nil {
		log.Fatal("❌ อัปเกรด Database ไม่สำเร็จ: ", err)
	}

	DB = database
	log.Println("✅ เชื่อมต่อ Database สำเร็จ")

}