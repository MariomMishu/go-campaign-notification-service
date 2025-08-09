package domain

import (
	"ems/models"
	"ems/types"
)

type (
	UserService interface {
		CreateUser(req *types.CreateUserReq) error
		ReadUserByEmail(email string) (*models.User, error)
		StoreInCache(user *types.UserInfo) error
		ReadUser(id int) (*models.User, error)
		ReadPermissionsByRole(roleID int) ([]*models.Permission, error)
	}
	UserRepository interface {
		CreateUser(user *models.User) (*models.User, error)
		ReadUserById(id int) (*models.User, error)
		UserCountByEmail(email string) (int, error)
		ReadUserByEmail(email string) (*models.User, error)
		ReadPermissionsByRole(roleID int) ([]*models.Permission, error)
		ReadUsers(roleIds []int) ([]*models.User, error)
	}
)
