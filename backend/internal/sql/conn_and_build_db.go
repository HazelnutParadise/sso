package sql

import (
	"log"
	"sso/internal/sql/models"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

// 使用 GORM 連接資料庫並初始化
func ConnectAndBuildDB() {
	dsn := "host=localhost port=5432 user=你的帳號 password=你的密碼 dbname=你的資料庫 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("GORM 連線失敗: %v", err)
	}
	db.AutoMigrate(
		&models.User{},
		&models.UserProvider{},
		&models.OAuthClient{},
		&models.OAuthToken{},
		&models.LoginLog{},
	)
}
