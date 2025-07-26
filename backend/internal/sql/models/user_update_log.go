package models

import (
	"time"
)

// 使用者資訊更改紀錄
type UserUpdateLog struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`    // 外鍵
	Field     string    `gorm:"size:100;not null"` // 被更改的欄位
	OldValue  string    `gorm:"size:500"`
	NewValue  string    `gorm:"size:500"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserID"`
}
