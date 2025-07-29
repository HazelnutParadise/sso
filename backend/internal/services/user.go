package services

import (
	"errors"
	"sso/internal/logger"
	"sso/internal/sql"
	"sso/internal/sql/models"
	"time"

	"sso/internal/services/dto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 管理使用者資料
type userService struct{}

var (
	LoginErr_InvalidCredentials = errors.New("帳號或密碼錯誤")
	LoginErr_AccountSuspended   = errors.New("此帳號已停權")
	LoginErr_PasswordNotSet     = errors.New("此帳號未設定密碼")
	LoginErr_UnknownError       = errors.New("登入時發生錯誤")
)

// 登入：驗證帳號密碼，成功回傳 user
func (s *userService) Login(email, password string) (*dto.UserDTO, error) {
	ip := getClientIP()    // 假設有方法獲取客戶端 IP
	userAgent := "unknown" // todo

	user, err := sql.GetUserByEmail(nil, email)
	if err != nil {
		errorType := LoginErr_InvalidCredentials
		if logErr := sql.AddLoginLog(nil,
			0,
			models.LoginMethodPassword, false,
			&ip, &userAgent,
			false, &errorType); logErr != nil {
			logger.Log.WithError(logErr).Error("登入失敗，記錄日誌失敗")
		}
		return nil, errorType
	}
	if user.IsActive == false {
		errorType := LoginErr_AccountSuspended
		if logErr := sql.AddLoginLog(nil,
			0,
			models.LoginMethodPassword, false,
			&ip, &userAgent,
			false, &errorType); logErr != nil {
			logger.Log.WithError(logErr).Error("登入失敗，記錄日誌失敗")
		}
		return nil, errorType
	}
	if user.PasswordHash == nil {
		// 如果沒有設定密碼，則無法登入
		errorType := LoginErr_PasswordNotSet
		if logErr := sql.AddLoginLog(nil,
			0,
			models.LoginMethodPassword, false,
			&ip, &userAgent,
			false, &errorType); logErr != nil {
			logger.Log.WithError(logErr).Error("登入失敗，記錄日誌失敗")
		}
		return nil, errorType
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		errorType := LoginErr_InvalidCredentials
		if logErr := sql.AddLoginLog(nil,
			0,
			models.LoginMethodPassword, false,
			&ip, &userAgent,
			false, &errorType); logErr != nil {
			logger.Log.WithError(logErr).Error("登入失敗，記錄日誌失敗")
		}
		return nil, errorType
	}
	now := time.Now()
	user.LastLoginAt = &now

	err = sql.Transaction(func(tx *gorm.DB) error {
		// 更新最後登入時間
		if err := sql.UpdateUser(tx, user); err != nil {
			return err
		}
		if err := sql.AddLoginLog(tx,
			user.ID,
			models.LoginMethodPassword, false,
			&ip, &userAgent,
			true, nil, // 登入成功
		); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, LoginErr_UnknownError
	}
	return dto.ToUserDTO(user), nil
}

// 登出：可記錄登出日誌
func (s *userService) Logout(userID uint) error {
	// 實際可記錄登出日誌或做其它處理
	// todo
	return nil
}

func getClientIP() string {
	// 這裡應該從請求上下文中獲取 IP 地址
	// 假設有 utils.GetClientIP() 方法
	// todo
	return ""
}
