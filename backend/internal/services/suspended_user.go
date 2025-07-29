package services

import (
	"sso/internal/sql"
	"sso/internal/sql/models"
)

// SuspendedUserService 管理停權/解鎖帳號
type suspendedUserService struct{}

func (s *suspendedUserService) SuspendUser(userID uint, reason string, suspendedBy *uint) error {
	return sql.SuspendUser(nil, userID, reason, suspendedBy)
}

func (s *suspendedUserService) UnsuspendUser(userID uint) error {
	// 解鎖：將 is_active 改回 true
	user, err := sql.GetUserByID(nil, userID)
	if err != nil {
		return err
	}
	user.IsActive = true
	return sql.UpdateUser(nil, user)
}

func (s *suspendedUserService) GetSuspendedLogs(userID uint, limit int) ([]models.SuspendedUserLog, error) {
	logs, err := sql.GetSingleUserSuspendedLogs(nil, userID, limit)
	if err != nil {
		return nil, err
	}
	// 假設有 dto.ToSuspendedUserLogDTO
	// return dto.ModelToDTO(&logs[0], dto.ToSuspendedUserLogDTO), nil
	return logs, nil // TODO: 實作 DTO 轉換
}
