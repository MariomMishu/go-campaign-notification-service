package domain

import (
	"ems/types"
)

type (
	AuthService interface {
		Login(req *types.LoginReq) (*types.LoginResp, error)
		VerifyAccessToken(tokenString string) (*types.UserInfo, *types.Token, error)
		Logout(accessTokenUuid, refreshTokenUuid string) error
	}
)
