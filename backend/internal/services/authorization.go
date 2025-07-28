package services

import (
	"errors"
	"sso/internal/sql"
)

// AuthorizationService 負責權限判斷
type authorizationService struct{}

func (s *authorizationService) Authorize(userID uint, resource string, action string) (bool, error) {
	user, err := sql.GetUserByID(userID)
	if err != nil {
		return false, errors.New("找不到使用者")
	}
	if !user.IsActive {
		return false, errors.New("使用者未啟用")
	}
	// TODO: 根據 resource/action 實作更細緻的權限
	return true, nil
}
