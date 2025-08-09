package services

import (
	"ems/consts"
	"ems/domain"
	"ems/types"
	"fmt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userSvc  domain.UserService
	tokenSvc domain.TokenService
}

func NewAuthServiceImpl(userSvc domain.UserService, tokenSvc domain.TokenService) *AuthServiceImpl {
	return &AuthServiceImpl{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
	}
}
func (authSvc *AuthServiceImpl) Login(req *types.LoginReq) (*types.LoginResp, error) {

	user, err := authSvc.userSvc.ReadUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err

	}

	token, err := authSvc.tokenSvc.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	if err := authSvc.tokenSvc.StoreTokenUUID(token); err != nil {
		fmt.Println("error occurred while storing token uuid", err)
		return nil, err
	}
	userInfo := &types.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RoleID:    user.RoleID,
		Role:      consts.RoleMap[user.RoleID],
	}

	resp := &types.LoginResp{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		User:         userInfo,
	}
	go func() {
		if err := authSvc.userSvc.StoreInCache(userInfo); err != nil {
			fmt.Println("error occurred while user store in cache", err)
		}
	}()
	return resp, nil
}
func (authSvc *AuthServiceImpl) VerifyAccessToken(tokenString string) (*types.UserInfo, *types.Token, error) {
	token, err := authSvc.tokenSvc.ParseAccessToken(tokenString)
	if err != nil {
		return nil, nil, err
	}

	cachedUserID, err := authSvc.tokenSvc.ReadUserIDFromAccessTokenUUID(token.AccessUuid)

	if err != nil {
		fmt.Println("error occurred while reading user from cache", err)
		return nil, nil, err
	}

	if cachedUserID != token.UserID {
		return nil, nil, fmt.Errorf("user not found")
	}
	user, err := authSvc.userSvc.ReadUser(cachedUserID)
	if err != nil {
		fmt.Println("error occurred while reading user by Id", err)
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, fmt.Errorf("user not found")

	}
	userInfo := &types.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RoleID:    user.RoleID,
		Role:      consts.RoleMap[user.RoleID],
	}

	return userInfo, token, nil

}
func (authSvc *AuthServiceImpl) Logout(accessTokenUuid, refreshTokenUuid string) error {
	var token = &types.Token{AccessUuid: accessTokenUuid, RefreshUuid: refreshTokenUuid}
	if err := authSvc.tokenSvc.DeleteTokenUUID(token); err != nil {
		log.Printf("Token deletion error: %v", err)
	}
	return nil
}
