package services

import (
	"sso/internal/sql"
	"sso/internal/sql/models"
)

// LogService 負責記錄系統事件
type logService struct{}

func (s *logService) GetUserLoginLogs(userID uint, limit int) ([]models.LoginLog, error) {
	logs, err := sql.GetUserLoginLogs(nil, userID, limit)
	if err != nil {
		return nil, err
	}
	// 假設有 dto.ToLoginLogDTO
	// return dto.ModelToDTO(&logs[0], dto.ToLoginLogDTO), nil
	return logs, nil // TODO: 實作 DTO 轉換
}
