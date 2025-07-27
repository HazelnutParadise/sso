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

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func AddUserProvider(provider *models.UserProvider) error {
	return db.Create(provider).Error
}

func GetUserProviders(userID uint) ([]models.UserProvider, error) {
	var providers []models.UserProvider
	if err := db.Where("user_id = ?", userID).Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

func DeleteUserProvider(userID uint, providerName string) error {
	return db.Where("user_id = ? AND provider_name = ?", userID, providerName).Delete(&models.UserProvider{}).Error
}

func GetUserByProviderUserID(providerName, providerUserID string) (*models.User, error) {
	var provider models.UserProvider
	if err := db.Where("provider_name = ? AND provider_user_id = ?", providerName, providerUserID).First(&provider).Error; err != nil {
		return nil, err
	}

	var user models.User
	if err := db.First(&user, provider.UserID).Error; err != nil {
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

func GetUserUpdateLogs(userID uint, limit int) ([]models.UserUpdateLog, error) {
	// 如果 limit 為 -1，則表示不限制數量
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.UserUpdateLog
	query := db.Where("user_id = ?", userID).Order("updated_at desc").Limit(limit)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func GetUserPasswordUpdateLogs(userID uint, limit int) ([]models.UserUpdateLog, error) {
	// 如果 limit 為 -1，則表示不限制數量
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.UserUpdateLog
	query := db.Where("user_id = ? AND field = ?", userID, "password_hash").Order("updated_at desc").Limit(limit)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
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

func DeleteUser(userID uint) error {
	return db.Delete(&models.User{}, userID).Error
}

func AddOauthClient(client *models.OAuthClient) error {
	return db.Create(client).Error
}

func GetOauthClientByID(clientID uint) (*models.OAuthClient, error) {
	var client models.OAuthClient
	if err := db.First(&client, clientID).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func GetOauthClientByClientID(clientID string) (*models.OAuthClient, error) {
	var client models.OAuthClient
	if err := db.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func UpdateOauthClient(client *models.OAuthClient) error {
	// 先查詢原始資料
	var oldClient models.OAuthClient
	if err := db.First(&oldClient, client.ID).Error; err != nil {
		return err
	}
	return db.Save(client).Error
}

// 取得 string pointer 的值，若為 nil 則回傳空字串
func getStringPointerValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
