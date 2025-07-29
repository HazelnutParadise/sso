package services

import (
	"errors"
	"sso/internal/sql"
	"sso/internal/sql/models"
	"time"

	"sso/internal/services/dto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 管理使用者資料
type userService struct{}

// 登入：驗證帳號密碼，成功回傳 user
func (s *userService) Login(email, password string) (*dto.UserDTO, error) {
	user, err := sql.GetUserByEmail(nil, email)
	if err != nil {
		return nil, errors.New("帳號或密碼錯誤")
	}
	if user.PasswordHash == nil {
		return nil, errors.New("此帳號未設定密碼")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("帳號或密碼錯誤")
	}
	now := time.Now()
	user.LastLoginAt = &now

	ip := getClientIP() // 假設有方法獲取客戶端 IP

	userAgent := "unknown" // todo

	err = sql.Transaction(func(tx *gorm.DB) error {
		// 更新最後登入時間
		if err := sql.UpdateUser(tx, user); err != nil {
			return err
		}
		if err := sql.AddLoginLog(tx,
			user.ID,
			models.LoginMethodPassword, false,
			&ip, &userAgent,
		); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("登入失敗，請稍後再試")
	}
	return dto.ModelToDTO(user, dto.ToUserDTO), nil
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
