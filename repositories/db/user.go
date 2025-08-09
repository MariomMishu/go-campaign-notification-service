package db

import (
	"ems/models"
)

func (repo *Repository) CreateUser(user *models.User) (*models.User, error) {
	err := repo.db.Create(user).Error
	return user, err
}

func (repo *Repository) UserCountByEmail(email string) (int, error) {
	var count int64
	if err := repo.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
func (repo *Repository) ReadUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := repo.db.Model(&models.User{}).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}
func (repo *Repository) ReadUserById(id int) (*models.User, error) {
	var user models.User
	if err := repo.db.Model(&models.User{}).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}

func (repo *Repository) ReadPermissionsByRole(roleID int) ([]*models.Permission, error) {
	var permissions []*models.Permission

	if err := repo.db.Model(&models.RolePermission{}).
		Select("permissions.*").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

func (repo *Repository) ReadUsers(roleIds []int) ([]*models.User, error) {
	var users []*models.User
	if err := repo.db.Model(&models.User{}).Where("role_id IN (?)", roleIds).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
