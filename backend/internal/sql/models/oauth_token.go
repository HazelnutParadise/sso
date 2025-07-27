package models

import "time"

// 🔹 OAuth Token 表（記錄授權 token）
type OAuthToken struct {
	ID           uint    `gorm:"primaryKey"`
	UserID       uint    `gorm:"not null"`
	ClientID     uint    `gorm:"not null"`
	AccessToken  string  `gorm:"size:500;not null"`
	RefreshToken *string `gorm:"size:500"`
	ExpiresAt    *time.Time
	Scope        string     `gorm:"size:500"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `gorm:"index"` // 軟刪除欄位

	// 外鍵
	User        User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OAuthClient OAuthClient `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
