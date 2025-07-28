package dto

import "sso/internal/sql/models"

type OAuthClientDTO struct {
	ID          uint   `json:"id"`
	ClientID    string `json:"client_id"`
	Name        string `json:"name"`
	RedirectURI string `json:"redirect_uri"`
	Scopes      string `json:"scopes"`
}

func ToOAuthClientDTO(c *models.OAuthClient) *OAuthClientDTO {
	if c == nil {
		return nil
	}
	return &OAuthClientDTO{
		ID:          c.ID,
		ClientID:    c.ClientID,
		Name:        c.Name,
		RedirectURI: c.RedirectURI,
		Scopes:      c.Scopes,
	}
}
