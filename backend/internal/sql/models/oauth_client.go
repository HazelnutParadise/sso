package models

import "time"

// 🔹 OAuth Client 表（提供 OAuth 給別人）
type OAuthClient struct {
	ID           uint   `gorm:"primaryKey"`
	ClientID     string `gorm:"uniqueIndex;size:100;not null"`
	ClientSecret string `gorm:"size:255;not null"`
	Name         string `gorm:"size:255;not null"`
	RedirectURI  string `gorm:"size:500;not null"`
	Scopes       string `gorm:"size:500"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// 關聯
	OAuthTokens []OAuthToken `gorm:"foreignKey:ClientID"`
}
