package services

// TokenService 管理 access/refresh token
type TokenService struct{}

func (s *TokenService) GenerateToken(userID string) (string, error) {
	// 產生 token
	return "token", nil
}

func (s *TokenService) ValidateToken(token string) (bool, error) {
	// 驗證 token
	return true, nil
}
