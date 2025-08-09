package domain

import "ems/types"

type (
	TokenService interface {
		CreateToken(UserID int) (*types.Token, error)
		StoreTokenUUID(token *types.Token) error
		ParseAccessToken(accessToken string) (*types.Token, error)
		ReadUserIDFromAccessTokenUUID(accessTokenUuid string) (int, error)
		DeleteTokenUUID(token *types.Token) error
	}
)
