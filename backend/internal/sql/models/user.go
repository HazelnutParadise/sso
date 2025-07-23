package models

import (
	"time"
)

// 🔹 使用者主表
type User struct {
	ID           uint    `gorm:"primaryKey"`
	Email        string  `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash *string `gorm:"size:255"` // 社群登入時可以是 NULL
	Name         *string `gorm:"size:100"`
	AvatarURL    *string `gorm:"size:500"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time
	IsActive     bool `gorm:"default:true"`

	// 關聯
	UserProviders []UserProvider `gorm:"foreignKey:UserID"`
	LoginLogs     []LoginLog     `gorm:"foreignKey:UserID"`
	OAuthTokens   []OAuthToken   `gorm:"foreignKey:UserID"`
}
