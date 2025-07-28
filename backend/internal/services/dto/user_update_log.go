package dto

import (
	"sso/internal/sql/models"
	"time"
)

type UserUpdateLogDTO struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Field     string `json:"field"`
	OldValue  string `json:"old_value"`
	NewValue  string `json:"new_value"`
	UpdatedAt string `json:"updated_at"`
}

func ToUserUpdateLogDTO(log *models.UserUpdateLog) *UserUpdateLogDTO {
	if log == nil {
		return nil
	}
	return &UserUpdateLogDTO{
		ID:        log.ID,
		UserID:    log.UserID,
		Field:     log.Field,
		OldValue:  log.OldValue,
		NewValue:  log.NewValue,
		UpdatedAt: formatTime(&log.UpdatedAt),
	}
}

func ToUserUpdateLogDTOs(logs []models.UserUpdateLog) []UserUpdateLogDTO {
	dtos := make([]UserUpdateLogDTO, 0, len(logs))
	for _, l := range logs {
		dto := ToUserUpdateLogDTO(&l)
		if dto != nil {
			dtos = append(dtos, *dto)
		}
	}
	return dtos
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
