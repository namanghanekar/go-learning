package config

import (
	"fmt"
	"log"

	"lmsbackend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ DB connection failed:", err)
	}

	DB = db

	// Auto migrate
	DB.AutoMigrate(
		&models.User{},
		&models.CreatorProfile{},
		&models.BrandProfile{},
	)

	fmt.Println("✅ Database connected")
}
