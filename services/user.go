package services

import (
	"ems/config"
	"ems/domain"
	"ems/models"
	"ems/types"
	"ems/utils/errutil"
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	repo     domain.UserRepository
	redisSvc *RedisService
}

func NewUserServiceImpl(userRepo domain.UserRepository, redisSvc *RedisService) *UserServiceImpl {
	return &UserServiceImpl{
		repo:     userRepo,
		redisSvc: redisSvc,
	}
}
func (svc *UserServiceImpl) CreateUser(req *types.CreateUserReq) error {
	isExist, err := svc.IsEmailExist(req.Email)
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("email already exists")
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	user := &models.User{
		Email:     req.Email,
		Password:  string(hashedPass),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleID:    req.RoleId,
	}
	if _, err := svc.repo.CreateUser(user); err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

func (svc *UserServiceImpl) IsEmailExist(email string) (bool, error) {
	// Check if the email already exists in the repository
	count, err := svc.repo.UserCountByEmail(email)
	if err != nil {
		fmt.Println("error occurred while fetching user by email", err)

		return false, err
	}
	return count != 0, nil
}
func (svc *UserServiceImpl) ReadUserByEmail(email string) (*models.User, error) {
	// Fetch user by email from the repository
	user, err := svc.repo.ReadUserByEmail(email)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("Error occurred while fetching user by email: %s, error: %v\n", email, err)
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errutil.ErrRecordNotFound
	}
	return user, nil
}
func (svc *UserServiceImpl) StoreInCache(user *types.UserInfo) error {
	userCacheKey := config.Redis().MandatoryPrefix + config.Redis().UserPrefix + strconv.Itoa(user.ID)

	if err := svc.redisSvc.SetStruct(userCacheKey, user, config.Redis().UserCacheTTL); err != nil {
		return fmt.Errorf("could not cache user in redis, id: [%d], err: [%v]", user.ID, err)
	}
	return nil
}
func (svc *UserServiceImpl) ReadUser(id int) (*models.User, error) {
	return svc.repo.ReadUserById(id)
}
func (svc *UserServiceImpl) ReadPermissionsByRole(roleID int) ([]*models.Permission, error) {
	return svc.repo.ReadPermissionsByRole(roleID)
}

// Define methods for user service
