package services

// SuspendedUserService 管理停權/解鎖帳號
type SuspendedUserService struct{}

func (s *SuspendedUserService) SuspendUser(userID string) error {
	// 停權邏輯
	return nil
}

func (s *SuspendedUserService) UnsuspendUser(userID string) error {
	// 解鎖邏輯
	return nil
}
