package services

import (
	"errors"
	"sso/internal/sql"
	"sso/internal/sql/models"
	"sso/internal/utils"
	"time"

	"sso/internal/services/dto"

	"golang.org/x/crypto/bcrypt"
)

// UserService 管理使用者資料
type userService struct{}

// 登入：驗證帳號密碼，成功回傳 user
func (s *userService) Login(email, password string) (*dto.UserDTO, error) {
	user, err := sql.GetUserByEmail(email)
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
	_ = sql.UpdateUser(user)
	return dto.ModelToDTO(user, dto.ToUserDTO), nil
}

// 登出：可記錄登出日誌
func (s *userService) Logout(userID uint) error {
	// 實際可記錄登出日誌或做其它處理
	// todo
	return nil
}

// 取得使用者
func (s *userService) GetUser(userID uint) (*dto.UserDTO, error) {
	user, err := sql.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return dto.ModelToDTO(user, dto.ToUserDTO), nil
}

// 依 email 查詢
func (s *userService) GetUserByEmail(email string) (*dto.UserDTO, error) {
	user, err := sql.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return dto.ModelToDTO(user, dto.ToUserDTO), nil
}

// 建立使用者
func (s *userService) CreateUser(user *models.User, password string) (*dto.UserDTO, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = utils.PtrString(string(hash))
	if err := sql.AddUser(user); err != nil {
		return nil, err
	}
	return dto.ModelToDTO(user, dto.ToUserDTO), nil
}

// 更新使用者
func (s *userService) UpdateUser(user *models.User) (*dto.UserDTO, error) {
	if err := sql.UpdateUser(user); err != nil {
		return nil, err
	}
	return dto.ModelToDTO(user, dto.ToUserDTO), nil
}

// 刪除使用者
func (s *userService) DeleteUser(userID uint) (*dto.UserDTO, error) {
	user, err := sql.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if err := sql.DeleteUser(userID); err != nil {
		return nil, err
	}
	return dto.ModelToDTO(user, dto.ToUserDTO), nil
}

// 取得異動紀錄
func (s *userService) GetUserUpdateLogs(userID uint, limit int) ([]dto.UserUpdateLogDTO, error) {
	logs, err := sql.GetUserUpdateLogs(userID, limit)
	if err != nil {
		return nil, err
	}
	return dto.ToUserUpdateLogDTOs(logs), nil
}

// 取得密碼異動紀錄
func (s *userService) GetUserPasswordUpdateLogs(userID uint, limit int) ([]dto.UserUpdateLogDTO, error) {
	logs, err := sql.GetUserPasswordUpdateLogs(userID, limit)
	if err != nil {
		return nil, err
	}
	return dto.ToUserUpdateLogDTOs(logs), nil
}
