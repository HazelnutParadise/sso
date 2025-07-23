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

func UpdateUser(user *models.User) error {
	toUpdateFields := []string{}
	if user.Name != nil && *user.Name != "" {
		toUpdateFields = append(toUpdateFields, "name")
	}
	if user.AvatarURL != nil && *user.AvatarURL != "" {
		toUpdateFields = append(toUpdateFields, "avatar_url")
	}
	if user.Email != "" {
		toUpdateFields = append(toUpdateFields, "email")
	}
	if user.PasswordHash != nil && *user.PasswordHash != "" {
		toUpdateFields = append(toUpdateFields, "password_hash")
	}
	if len(toUpdateFields) == 0 {
		return nil // Nothing to update
	}
	// 使用 Select 方法指定要更新的欄位
	return db.Model(&models.User{}).Where("id = ?", user.ID).Select(toUpdateFields).Updates(user).Error
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

// Set limit to -1 to get all logs
func GetSuspendedUserLogs(userID uint, limit int) ([]models.SuspendedUserLog, error) {
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.SuspendedUserLog
	if err := db.Where("user_id = ?", userID).Order("suspended_at desc").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func GetLatestSuspendedUserLog(userID uint) (*models.SuspendedUserLog, error) {
	var log models.SuspendedUserLog
	if err := db.Where("user_id = ?", userID).Order("suspended_at desc").First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}
