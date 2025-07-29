package models

import (
	"time"

	"gorm.io/gorm"
)

// 🔹 使用者主表
type User struct {
	ID           uint           `gorm:"primaryKey"`
	Email        string         `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash *string        `gorm:"size:255"` // 社群登入時可以是 NULL
	Name         *string        `gorm:"size:100"`
	AvatarURL    *string        `gorm:"size:500"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	LastLoginAt  *time.Time     `gorm:"index"` // 最後登入時間
	IsActive     bool           `gorm:"default:true"`
	DeletedAt    gorm.DeletedAt `gorm:"index"` // 軟刪除欄位

	// 關聯
	UserProviders     []UserProvider     `gorm:"foreignKey:UserID"`
	LoginLogs         []LoginLog         `gorm:"foreignKey:UserID"`
	OAuthTokens       []OAuthToken       `gorm:"foreignKey:UserID"`
	SuspendedUserLogs []SuspendedUserLog `gorm:"foreignKey:UserID"`
	UserUpdateLogs    []UserUpdateLog    `gorm:"foreignKey:UserID"` // 使用者資訊更改紀錄
}
