package main

import (
	"go-rest-setup/src/database/models"
	config "go-rest-setup/src/lib/configs"
	"log"

	"gorm.io/gorm"
)

func main() {
	db := config.InitDatabase()

	if err := db.AutoMigrate(&models.User{}, &models.AuditLog{}, &models.File{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	createIndexSoftDelete(db)

}

func createIndexSoftDelete(db *gorm.DB) {
	if err := db.AutoMigrate(&models.User{}); err != nil {
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
