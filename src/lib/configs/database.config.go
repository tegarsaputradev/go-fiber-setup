package config

import (
	"fmt"
	"go-rest-setup/src/database/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func createIndexSoftDelete(db *gorm.DB) {
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	db.Exec(`
        CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique
        ON users(email)
        WHERE deleted_at IS NULL;
    `)

	db.Exec(`
        CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_unique
        ON users(username)
        WHERE deleted_at IS NULL;
    `)

	log.Println("Create Index Completed")
}

func InitDatabase() {
	c := EnvModule()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Pass,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	DB = db

	createIndexSoftDelete(db)
}
