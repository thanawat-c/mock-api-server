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

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disables prepared statements (Essential for Supabase Transaction Pooler 6543)
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            false,
	})
	if err != nil {
		log.Fatal("❌ เชื่อมต่อ Database ไม่สำเร็จ: ", err)
	}

	// Optimize Connection Pool (ท่อส่งข้อมูลค้างไว้ ไม่ต้องต่อใหม่ทุกครั้ง)
	sqlDB, err := database.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
	}

	if err := database.AutoMigrate(&models.Endpoint{}, &models.MockRecord{}); err != nil {
		log.Println("⚠️ AutoMigrate warning: ", err)
	}

	DB = database
	log.Println("✅ เชื่อมต่อ Database สำเร็จ พร้อมใช้งาน Connection Pool")

}