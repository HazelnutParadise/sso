package sql

import (
	"sso/internal/sql/models"
	"time"

	"gorm.io/gorm"
)

var Transaction = db.Transaction

// getDBOrTx: tx 為 nil 時自動用 db
func getDBOrTx(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return db
}

func AddUser(tx *gorm.DB, user *models.User) error {
	dbx := getDBOrTx(tx)
	return dbx.Create(user).Error
}

func GetUserByID(tx *gorm.DB, userID uint) (*models.User, error) {
	dbx := getDBOrTx(tx)
	var user models.User
	if err := dbx.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(tx *gorm.DB, email string) (*models.User, error) {
	dbx := getDBOrTx(tx)
	var user models.User
	if err := dbx.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func AddUserProvider(tx *gorm.DB, provider *models.UserProvider) error {
	dbx := getDBOrTx(tx)
	return dbx.Create(provider).Error
}

func GetUserProviders(tx *gorm.DB, userID uint) ([]models.UserProvider, error) {
	dbx := getDBOrTx(tx)
	var providers []models.UserProvider
	if err := dbx.Where("user_id = ?", userID).Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

func DeleteUserProvider(tx *gorm.DB, userID uint, providerName string) error {
	dbx := getDBOrTx(tx)
	return dbx.Where("user_id = ? AND provider_name = ?", userID, providerName).Delete(&models.UserProvider{}).Error
}

func GetUserByProviderUserID(tx *gorm.DB, providerName, providerUserID string) (*models.User, error) {
	dbx := getDBOrTx(tx)
	var provider models.UserProvider
	if err := dbx.Where("provider_name = ? AND provider_user_id = ?", providerName, providerUserID).First(&provider).Error; err != nil {
		return nil, err
	}

	var user models.User
	if err := dbx.First(&user, provider.UserID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(tx *gorm.DB, user *models.User) error {
	dbx := getDBOrTx(tx)
	// 先查詢原始資料
	var oldUser models.User
	if err := dbx.First(&oldUser, user.ID).Error; err != nil {
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

	if err := dbx.Model(&models.User{}).Where("id = ?", user.ID).Select(toUpdateFields).Updates(user).Error; err != nil {
		return err
	}
	for _, log := range updateLogs {
		if err := dbx.Create(&log).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetUserUpdateLogs(tx *gorm.DB, userID uint, limit int) ([]models.UserUpdateLog, error) {
	// 如果 limit 為 -1，則表示不限制數量
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.UserUpdateLog
	query := tx.Where("user_id = ?", userID).Order("updated_at desc").Limit(limit)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func GetUserPasswordUpdateLogs(tx *gorm.DB, userID uint, limit int) ([]models.UserUpdateLog, error) {
	// 如果 limit 為 -1，則表示不限制數量
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.UserUpdateLog
	query := tx.Where("user_id = ? AND field = ?", userID, "password_hash").Order("updated_at desc").Limit(limit)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// SuspendUser 停權使用者並記錄原因
func SuspendUser(tx *gorm.DB, userID uint, reason string, suspendedBy *uint) error {
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
}

func GetSuspendedUsersLogs(tx *gorm.DB, limit int) ([]models.SuspendedUserLog, error) {
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.SuspendedUserLog
	query := tx.Order("suspended_at desc").Limit(limit)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// Set limit to -1 to get all logs
func GetSingleUserSuspendedLogs(tx *gorm.DB, userID uint, limit int) ([]models.SuspendedUserLog, error) {
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.SuspendedUserLog
	if err := tx.Where("user_id = ?", userID).Order("suspended_at desc").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func DeleteUser(tx *gorm.DB, userID uint) error {
	return tx.Delete(&models.User{}, userID).Error
}

func AddLoginLog(tx *gorm.DB,
	userID uint,
	loginMethod string, isOAuth bool,
	ipAddress *string, userAgent *string,
) error {
	log := &models.LoginLog{
		UserID:      userID,
		LoginMethod: loginMethod,
		IsOAuth:     isOAuth,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}
	dbx := getDBOrTx(tx)
	return dbx.Create(log).Error
}

func GetUserLoginLogs(tx *gorm.DB, userID uint, limit int) ([]models.LoginLog, error) {
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.LoginLog
	query := tx.Where("user_id = ?", userID).Order("login_at desc").Limit(limit)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func GetLoginLogsByIP(tx *gorm.DB, ipAddress string, limit int) ([]models.LoginLog, error) {
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.LoginLog
	query := tx.Where("ip_address = ?", ipAddress).Order("login_at desc").Limit(limit)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func GetUserLoginLogsByIP(tx *gorm.DB, userID uint, ipAddress string, limit int) ([]models.LoginLog, error) {
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.LoginLog
	query := tx.Where("user_id = ? AND ip_address = ?", userID, ipAddress).Order("login_at desc").Limit(limit)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func GetLoginLogsBetweenDates(tx *gorm.DB, start, end time.Time, limit int) ([]models.LoginLog, error) {
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.LoginLog
	query := tx.Where("login_at BETWEEN ? AND ?", start, end).Order("login_at desc").Limit(limit)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func GetUserLoginLogsBetweenDates(tx *gorm.DB, userID uint, start, end time.Time, limit int) ([]models.LoginLog, error) {
	if limit < -1 || limit == 0 {
		return nil, gorm.ErrInvalidValue
	}

	var logs []models.LoginLog
	query := tx.Where("user_id = ? AND (login_at BETWEEN ? AND ?)", userID, start, end).Order("login_at desc").Limit(limit)
	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func AddOauthClient(tx *gorm.DB, client *models.OAuthClient) error {
	return tx.Create(client).Error
}

func GetOauthClientByID(tx *gorm.DB, ID uint) (*models.OAuthClient, error) {
	var client models.OAuthClient
	if err := tx.First(&client, ID).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func GetOauthClientByClientID(tx *gorm.DB, clientID string) (*models.OAuthClient, error) {
	var client models.OAuthClient
	if err := tx.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func UpdateOauthClient(tx *gorm.DB, client *models.OAuthClient) error {
	// 先查詢原始資料
	var oldClient models.OAuthClient
	if err := tx.First(&oldClient, client.ID).Error; err != nil {
		return err
	}
	return tx.Save(client).Error
}

func DeleteOauthClient(tx *gorm.DB, ID uint) error {
	return tx.Delete(&models.OAuthClient{}, ID).Error
}

func AddOauthToken(tx *gorm.DB, token *models.OAuthToken) error {
	return tx.Create(token).Error
}

func GetOauthTokenByID(tx *gorm.DB, ID uint) (*models.OAuthToken, error) {
	var token models.OAuthToken
	if err := tx.First(&token, ID).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func GetOauthTokensByUserID(tx *gorm.DB, userID uint) ([]models.OAuthToken, error) {
	var tokens []models.OAuthToken
	if err := tx.Where("user_id = ?", userID).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func GetOauthTokensByClientID(tx *gorm.DB, clientID uint) ([]models.OAuthToken, error) {
	var tokens []models.OAuthToken
	if err := tx.Where("client_id = ?", clientID).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

// 依 access_token 查詢 OAuthToken
func GetOauthTokenByAccessToken(tx *gorm.DB, accessToken string) (*models.OAuthToken, error) {
	var token models.OAuthToken
	err := tx.Where("access_token = ?", accessToken).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func UpdateOauthToken(tx *gorm.DB, token *models.OAuthToken) error {
	// 先查詢原始資料
	var oldToken models.OAuthToken
	if err := tx.First(&oldToken, token.ID).Error; err != nil {
		return err
	}
	return tx.Save(token).Error
}

func DeleteOauthToken(tx *gorm.DB, ID uint) error {
	return tx.Delete(&models.OAuthToken{}, ID).Error
}

// 取得 string pointer 的值，若為 nil 則回傳空字串
func getStringPointerValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
