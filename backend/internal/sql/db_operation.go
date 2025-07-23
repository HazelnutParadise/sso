package sql

import (
	"sso/internal/sql/models"
	"time"

	"gorm.io/gorm"
)

func AddUser(user *models.User) error {
	return db.Create(user).Error
}

func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// SuspendUser 停權使用者並記錄原因
func SuspendUser(userID uint, reason string, suspendedBy *uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 更新使用者狀態
		if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("is_active", false).Error; err != nil {
			return err
		}
		// 2. 寫入停權紀錄
		log := models.SuspendedUserLog{
			UserID:      userID,
			Reason:      reason,
			SuspendedAt: time.Now(),
			SuspendedBy: suspendedBy,
		}
		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		return nil
	})
}
