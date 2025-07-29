package models

import "time"

// 🔹 登入紀錄表
type LoginLog struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null"`
	LoginMethod  string    `gorm:"size:50;not null"`       // password / google / github…
	IsOAuth      bool      `gorm:"not null;default:false"` // 是否為 OAuth 登入
	IPAddress    *string   `gorm:"size:45"`                // 可存 IPv4/IPv6
	UserAgent    *string   `gorm:"size:500"`
	IsSuccess    bool      `gorm:"not null;default:false"` // 是否登入成功
	ErrorMessage *string   `gorm:"size:500"`               // 登入失敗時的錯誤訊息
	AttemptedAt  time.Time `gorm:"autoCreateTime"`

	// 外鍵
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

const (
	LoginMethodPassword = "password"
	LoginMethodGoogle   = "google"
	LoginMethodGithub   = "github"
)
