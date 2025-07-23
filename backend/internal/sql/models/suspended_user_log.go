package models

import (
	"time"
)

// 停權紀錄表
// SuspendedUserLog 用來記錄使用者被停權的原因與時間
// UserID: 被停權的使用者ID
// Reason: 停權原因
// SuspendedAt: 停權時間
// SuspendedBy: 操作者ID（可選）
type SuspendedUserLog struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"index;not null"`
	Reason      string    `gorm:"size:255;not null"`
	SuspendedAt time.Time `gorm:"not null;autoCreateTime"`
	SuspendedBy *uint     `gorm:"index"`
}
