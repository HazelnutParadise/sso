package dto

import (
	"sso/internal/sql/models"
	"sso/internal/utils"
)

type LoginLogDTO struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	LoginMethod string `json:"login_method"`
	IsOAuth     bool   `json:"is_oauth"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
	LoginAt     string `json:"login_at"`
}

func ToLoginLogDTO(l *models.LoginLog) *LoginLogDTO {
	if l == nil {
		return nil
	}
	return &LoginLogDTO{
		ID:          l.ID,
		UserID:      l.UserID,
		LoginMethod: l.LoginMethod,
		IsOAuth:     l.IsOAuth,
		IPAddress:   utils.Deref(l.IPAddress),
		UserAgent:   utils.Deref(l.UserAgent),
		LoginAt:     utils.FormatTime(&l.LoginAt),
	}
}
