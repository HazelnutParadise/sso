package services

// LogService 負責記錄系統事件
type LogService struct{}

func (s *LogService) WriteLog(event string, userID string) error {
	// 寫入日誌
	return nil
}
