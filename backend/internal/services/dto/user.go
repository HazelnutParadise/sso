package dto

import (
	"sso/internal/sql/models"
	"sso/internal/utils"
)

type UserDTO struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name,omitempty"`
	AvatarBase64 string `json:"avatar_base64,omitempty"`
	IsActive     bool   `json:"is_active"`
	LastLoginAt  string `json:"last_login_at,omitempty"`
}

func ToUserDTO(u *models.User) *UserDTO {
	if u == nil {
		return nil
	}
	return &UserDTO{
		ID:           u.ID,
		Email:        u.Email,
		Name:         utils.Deref(u.Name),
		AvatarBase64: utils.BlobToBase64(*u.Avatar),
		IsActive:     u.IsActive,
		LastLoginAt:  utils.FormatTime(u.LastLoginAt),
	}
}
