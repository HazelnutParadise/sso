package sql

import (
	"fmt"
	"log"
	"sso/internal/env"
	"sso/internal/sql/models"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var db *gorm.DB

// 使用 GORM 連接資料庫並初始化
func ConnectAndInitDB() {
	host := env.DB_HOST
	port := env.DB_PORT
	user := env.DB_USER
	password := env.DB_PASSWORD
	dbname := env.DB_NAME

	// 首先連接到 postgres 預設資料庫檢查並創建目標資料庫
	createDatabaseIfNotExists(host, port, user, password, dbname)

	// 連接到目標資料庫
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("GORM 連線失敗: %v", err)
	}

	// 自動遷移資料表
	err = db.AutoMigrate(
		&models.User{},
		&models.UserProvider{},
		&models.OAuthClient{},
		&models.OAuthToken{},
		&models.LoginLog{},
		&models.SuspendedUserLog{},
	)
	if err != nil {
		log.Fatalf("資料表遷移失敗: %v", err)
	}

	log.Println("資料庫連接和初始化成功")
}

// 檢查並創建資料庫
func createDatabaseIfNotExists(host, port, user, password, dbname string) {
	// 連接到 postgres 預設資料庫
	defaultDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password)

	defaultDB, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("連接到 postgres 預設資料庫失敗: %v", err)
	}

	// 檢查資料庫是否存在
	var count int64
	sqlDB, err := defaultDB.DB()
	if err != nil {
		log.Fatalf("取得底層 SQL DB 連接失敗: %v", err)
	}

	err = sqlDB.QueryRow("SELECT 1 FROM pg_database WHERE datname = $1", dbname).Scan(&count)
	if err != nil {
		// 資料庫不存在，創建它
		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			log.Fatalf("創建資料庫失敗: %v", err)
		}
		log.Printf("資料庫 '%s' 創建成功", dbname)
	} else {
		log.Printf("資料庫 '%s' 已存在", dbname)
	}

	// 關閉預設資料庫連接
	sqlDB.Close()
}
