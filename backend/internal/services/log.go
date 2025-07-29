package services

import (
	"sso/internal/sql"
	"sso/internal/sql/models"
	"sso/internal/utils"
)

// LogService 負責記錄系統事件
type logService struct{}

func (s *logService) WriteLoginLog(userID uint, method, ip string, isOAuth bool, userAgent string) error {
	log := &models.LoginLog{
		UserID:      userID,
		LoginMethod: method,
		IsOAuth:     isOAuth,
		IPAddress:   utils.PtrString(ip),
		UserAgent:   utils.PtrString(userAgent),
	}
	return sql.AddLoginLog(log)
}

func (s *logService) GetUserLoginLogs(userID uint, limit int) ([]models.LoginLog, error) {
	logs, err := sql.GetUserLoginLogs(userID, limit)
	if err != nil {
		return nil, err
	}
	// 假設有 dto.ToLoginLogDTO
	// return dto.ModelToDTO(&logs[0], dto.ToLoginLogDTO), nil
	return logs, nil // TODO: 實作 DTO 轉換
}
