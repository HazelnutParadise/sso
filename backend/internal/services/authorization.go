package services

// AuthorizationService 負責權限判斷
type AuthorizationService struct{}

func (s *AuthorizationService) Authorize(userID string, resource string, action string) (bool, error) {
	// 授權邏輯
	return true, nil
}
