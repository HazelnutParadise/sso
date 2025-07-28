package services

// UserService 管理使用者資料
type UserService struct{}

func (s *UserService) GetUser(userID string) (interface{}, error) {
	// 查詢使用者
	return nil, nil
}

func (s *UserService) UpdateUser(userID string, data map[string]interface{}) error {
	// 更新使用者
	return nil
}
