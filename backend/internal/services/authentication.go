package services

import (
	"sso/internal/sql"

	"golang.org/x/crypto/bcrypt"
)

// AuthenticationService 驗證使用者身分
type authenticationService struct{}

func (s *authenticationService) Authenticate(email, password string) (bool, error) {
	user, err := sql.GetUserByEmail(email)
	if err != nil || user.PasswordHash == nil {
		return false, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return false, err
	}
	// 假設有 dto.ToUserDTO
	// dto.ModelToDTO(user, dto.ToUserDTO)
	return true, nil // TODO: 實作 DTO 轉換
}
