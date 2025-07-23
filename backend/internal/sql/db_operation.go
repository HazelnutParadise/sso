package sql

import "sso/internal/sql/models"

func AddUser(user *models.User) error {
	return db.Create(user).Error
}

func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func SuspendUser(userID uint) error {
	// todo: 紀錄停權原因(開一個新的表)
	return db.Model(&models.User{}).Where("id = ?", userID).Update("is_active", false).Error
}
