package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 🔹 社群綁定表
type UserProvider struct {
	ID              uint           `gorm:"primaryKey"`
	UserID          uint           `gorm:"not null"`
	ProviderName    string         `gorm:"size:50;not null"`
	ProviderUserID  string         `gorm:"size:255;not null"`
	ProviderEmail   *string        `gorm:"size:255"`
	ProviderRawJSON datatypes.JSON `gorm:"type:jsonb"`
	LinkedAt        time.Time      `gorm:"autoCreateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"` // 用於軟刪除

	// 外鍵
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
