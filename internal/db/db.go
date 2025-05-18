package db

import (
	"log"

	"github.com/Aadithya-J/code-sandbox/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	dsn := "postgresql://neondb_owner:npg_ILNzwcp8eGs4@ep-white-sky-a14ofyly-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated")
}
