package services

import (
	"errors"
	"regexp"
	"sso/internal/services/dto"
	"sso/internal/sql"
	"sso/internal/sql/models"
)

// OAuthClientService 管理第三方 client
type oauthClientService struct{}

// 註冊 client，clientID 必須唯一，redirectURI 必須為合法網址
func (s *oauthClientService) RegisterClient(client *models.OAuthClient) error {
	exist, _ := sql.GetOauthClientByClientID(nil, client.ClientID)
	if exist != nil {
		return errors.New("ClientID 已存在")
	}
	if !isValidURL(client.RedirectURI) {
		return errors.New("RedirectURI 格式錯誤")
	}
	return sql.AddOauthClient(nil, client)
}

func (s *oauthClientService) GetClientByID(id uint) (*models.OAuthClient, error) {
	client, err := sql.GetOauthClientByID(nil, id)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *oauthClientService) GetClientByClientID(clientID string) (*models.OAuthClient, error) {
	client, err := sql.GetOauthClientByClientID(nil, clientID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *oauthClientService) UpdateClient(client *models.OAuthClient) error {
	// 只能更新 name, redirectURI, scopes
	if !isValidURL(client.RedirectURI) {
		return errors.New("RedirectURI 格式錯誤")
	}
	return sql.UpdateOauthClient(nil, client)
}

func (s *oauthClientService) DeleteClient(id uint) error {
	return sql.DeleteOauthClient(nil, id)
}

// 取得 DTO
func (s *oauthClientService) GetClientDTOByID(id uint) (*dto.OAuthClientDTO, error) {
	client, err := s.GetClientByID(id)
	if err != nil {
		return nil, err
	}
	return dto.ToOAuthClientDTO(client), nil
}

func (s *oauthClientService) GetClientDTOByClientID(clientID string) (*dto.OAuthClientDTO, error) {
	client, err := s.GetClientByClientID(clientID)
	if err != nil {
		return nil, err
	}
	return dto.ToOAuthClientDTO(client), nil
}

// 工具：簡易 URL 格式檢查
func isValidURL(url string) bool {
	re := regexp.MustCompile(`^https?://[\w\-\.]+(:\d+)?(/.*)?$`)
	return re.MatchString(url)
}
