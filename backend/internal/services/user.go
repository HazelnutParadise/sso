package services

import (
	"errors"
	"sso/internal/logger"
	"sso/internal/sql"
	"sso/internal/sql/models"

	"sso/internal/services/dto"

	"sso/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 管理使用者資料
type userService struct{}

var (
	ErrLogin_InvalidCredentials = errors.New("帳號或密碼錯誤")
	ErrLogin_AccountSuspended   = errors.New("此帳號已停權")
	ErrLogin_PasswordNotSet     = errors.New("此帳號未設定密碼")
	ErrLogin_UnknownError       = errors.New("登入時發生錯誤")

	ErrRegister_EmailExists     = errors.New("此 email 已被註冊")
	ErrRegister_PasswordHashing = errors.New("密碼雜湊失敗")
)

// 登入：驗證帳號密碼，成功回傳 user
func (s *userService) Login(email, password string) (*dto.UserDTO, error) {
	ip := getClientIP()    // 假設有方法獲取客戶端 IP
	userAgent := "unknown" // todo
	user, err := loginCheckAndReturnUser(email, password)
	if err != nil {
		if logErr := sql.AddLoginLog(nil,
			0,
			models.LoginMethodPassword, false,
			&ip, &userAgent,
			false, &err); logErr != nil {
			logger.Log.WithError(logErr).Error("登入失敗，記錄日誌失敗")
		}
	}

	err = sql.Transaction(func(tx *gorm.DB) error {
		// 更新最後登入時間
		if err := sql.UpdateUserLoginTime(tx, user.ID); err != nil {
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
		logger.Log.WithError(err).Error("登入失敗，資料已回滾")
		return nil, ErrLogin_UnknownError
	}
	return user, nil
}

// 登出：可記錄登出日誌
func (s *userService) Logout(userID uint) error {
	// 實際可記錄登出日誌或做其它處理
	// todo
	return nil
}

func (s *userService) Register(
	email, password string, name string, avatarBase64 string,
) (*dto.UserDTO, error) {
	// 先檢查 email 是否已被註冊
	existingUser, err := sql.GetUserByEmail(nil, email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrRegister_EmailExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MaxCost)
	if err != nil {
		return nil, ErrRegister_PasswordHashing
	}
	passwordHashStr := string(passwordHash)

	// 創建新使用者
	newUser := &models.User{
		Email:        email,
		Name:         &name,
		Avatar:       utils.Base64ToBlob(avatarBase64),
		PasswordHash: &passwordHashStr,
	}
	if err := sql.AddUser(nil, newUser); err != nil {
		return nil, err
	}

	return dto.ToUserDTO(newUser), nil
}

func loginCheckAndReturnUser(email string, password string) (*dto.UserDTO, error) {
	user, err := sql.GetUserByEmail(nil, email)
	if err != nil {
		return nil, ErrLogin_InvalidCredentials
	}
	if !user.IsActive {
		return nil, ErrLogin_AccountSuspended
	}
	if user.PasswordHash == nil {
		// 如果沒有設定密碼，則無法登入
		return nil, ErrLogin_PasswordNotSet
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrLogin_InvalidCredentials
	}
	dtoUser := dto.ToUserDTO(user)
	return dtoUser, nil
}

func getClientIP() string {
	// 這裡應該從請求上下文中獲取 IP 地址
	// 假設有 utils.GetClientIP() 方法
	// todo
	return ""
}
