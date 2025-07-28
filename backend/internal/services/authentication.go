package services

// AuthenticationService 驗證使用者身分
type AuthenticationService struct{}

func (s *AuthenticationService) Authenticate(username, password string) (bool, error) {
	// 驗證邏輯
	return true, nil
}
