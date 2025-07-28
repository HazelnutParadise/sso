package services

import (
	"errors"
	"sso/internal/sql"
	"sso/internal/sql/models"
	"time"
)

// TokenService 管理 access/refresh token
type tokenService struct{}

// 產生 token
func (s *tokenService) GenerateToken(userID, clientID uint, scope string, expiresIn time.Duration) (*models.OAuthToken, error) {
	token := &models.OAuthToken{
		UserID:      userID,
		ClientID:    clientID,
		AccessToken: generateRandomToken(),
		Scope:       scope,
		ExpiresAt:   ptrTime(time.Now().Add(expiresIn)),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := sql.AddOauthToken(token)
	// 假設有 dto.ToTokenDTO
	// return dto.ModelToDTO(token, dto.ToTokenDTO), err
	return token, err // TODO: 實作 DTO 轉換
}

// 驗證 token
func (s *tokenService) ValidateToken(accessToken string) (*models.OAuthToken, error) {
	token, err := sql.GetOauthTokenByAccessToken(accessToken)
	if err != nil {
		return nil, errors.New("token 無效")
	}
	if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token 已過期")
	}
	return token, nil
}

// 失效 token
func (s *tokenService) RevokeToken(tokenID uint) error {
	return sql.DeleteOauthToken(tokenID)
}

// 工具
func generateRandomToken() string {
	// 這裡可用 crypto/rand 或其它安全隨機字串
	return "random-token" // TODO: 實作安全隨機
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
