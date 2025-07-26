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
	// 先查詢原始資料
	var oldUser models.User
	if err := db.First(&oldUser, user.ID).Error; err != nil {
		return err
	}

	toUpdateFields := []string{}
	updateLogs := []models.UserUpdateLog{}
	now := time.Now()

	if user.Name != nil && *user.Name != "" && (oldUser.Name == nil || *oldUser.Name != *user.Name) {
		toUpdateFields = append(toUpdateFields, "name")
		updateLogs = append(updateLogs, models.UserUpdateLog{
			UserID:    user.ID,
			Field:     "name",
			OldValue:  getStringPointerValue(oldUser.Name),
			NewValue:  *user.Name,
			UpdatedAt: now,
		})
	}
	if user.AvatarURL != nil && *user.AvatarURL != "" && (oldUser.AvatarURL == nil || *oldUser.AvatarURL != *user.AvatarURL) {
		toUpdateFields = append(toUpdateFields, "avatar_url")
		updateLogs = append(updateLogs, models.UserUpdateLog{
			UserID:    user.ID,
			Field:     "avatar_url",
			OldValue:  getStringPointerValue(oldUser.AvatarURL),
			NewValue:  *user.AvatarURL,
			UpdatedAt: now,
		})
	}
	if user.Email != "" && oldUser.Email != user.Email {
		toUpdateFields = append(toUpdateFields, "email")
		updateLogs = append(updateLogs, models.UserUpdateLog{
			UserID:    user.ID,
			Field:     "email",
			OldValue:  oldUser.Email,
			NewValue:  user.Email,
			UpdatedAt: now,
		})
	}
	if user.PasswordHash != nil && *user.PasswordHash != "" && (oldUser.PasswordHash == nil || *oldUser.PasswordHash != *user.PasswordHash) {
		toUpdateFields = append(toUpdateFields, "password_hash")
		updateLogs = append(updateLogs, models.UserUpdateLog{
			UserID:    user.ID,
			Field:     "password_hash",
			OldValue:  getStringPointerValue(oldUser.PasswordHash),
			NewValue:  *user.PasswordHash,
			UpdatedAt: now,
		})
	}
	if len(toUpdateFields) == 0 {
		return nil // Nothing to update
	}

	toUpdateFields = append(toUpdateFields, "updated_at") // Always update updated_at
	user.UpdatedAt = now

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Select(toUpdateFields).Updates(user).Error; err != nil {
			return err
		}
		for _, log := range updateLogs {
			if err := tx.Create(&log).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// 取得 string pointer 的值，若為 nil 則回傳空字串
func getStringPointerValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
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
