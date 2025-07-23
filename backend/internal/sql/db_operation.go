package sql

import "sso/internal/sql/models"

func AddUser(user *models.User) error {
	return db.Create(user).Error
}
