package models

import "time"

// 🔹 登入紀錄表
type LoginLog struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	LoginMethod string    `gorm:"size:50;not null"` // password / google / github…
	IPAddress   *string   `gorm:"size:45"`          // 可存 IPv4/IPv6
	UserAgent   *string   `gorm:"size:500"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	// 外鍵
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
