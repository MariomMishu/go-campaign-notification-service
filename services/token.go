package services

import (
	"ems/config"
	"ems/types"
	"ems/utils/errutil"
	"ems/utils/methodutil"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenServiceImpl struct {
	redisSvc *RedisService
}

func NewTokenServiceImpl(redisSvc *RedisService) *TokenServiceImpl {
	return &TokenServiceImpl{
		redisSvc: redisSvc,
	}
}

func (svc *TokenServiceImpl) CreateToken(userID int) (*types.Token, error) {

	jwtConf := config.Jwt()
	if jwtConf == nil {
		fmt.Println("JWT config not initialized for user ID:", userID)
		return nil, errors.New("JWT config not initialized")
	}
	token := &types.Token{}

	accessUuid := uuid.New().String()
	refreshUuid := uuid.New().String()

	token.UserID = userID
	token.AccessUuid = accessUuid
	token.RefreshUuid = refreshUuid
	token.AccessExpiry = time.Now().Add(time.Second * jwtConf.AccessTokenExpiry).Unix()
	token.RefreshExpiry = time.Now().Add(time.Second * jwtConf.RefreshTokenExpiry).Unix()
	//claims basically a map that holding the interface as value and key is the string
	atClaims := jwt.MapClaims{}
	atClaims["uid"] = userID
	atClaims["aid"] = token.AccessUuid
	atClaims["rid"] = token.RefreshUuid
	atClaims["exp"] = token.AccessExpiry
	var err error
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token.AccessToken, err = at.SignedString([]byte(jwtConf.AccessTokenSecret))
	if err != nil {
		return nil, errutil.ErrAccessTokenSign
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["uid"] = userID
	rtClaims["aid"] = token.AccessUuid
	rtClaims["rid"] = token.RefreshUuid
	rtClaims["exp"] = token.RefreshExpiry
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	token.RefreshToken, err = rt.SignedString([]byte(jwtConf.RefreshTokenSecret))

	if err != nil {
		return nil, errutil.ErrRefreshTokenSign
	}

	return token, nil
}
func (svc *TokenServiceImpl) StoreTokenUUID(token *types.Token) error {
	accessTokenCachekey := config.Redis().MandatoryPrefix + config.Redis().AccessUuidPrefix + token.AccessUuid
	err := svc.redisSvc.Set(accessTokenCachekey, token.UserID, time.Duration(token.AccessExpiry))

	if err != nil {
		return err
	}

	refreshTokenCachekey := config.Redis().MandatoryPrefix + config.Redis().RefreshUuidPrefix + token.RefreshUuid
	err = svc.redisSvc.Set(refreshTokenCachekey, token.UserID, time.Duration(token.RefreshExpiry))

	if err != nil {
		return err
	}
	return nil
}

func (svc *TokenServiceImpl) ParseAccessToken(accessToken string) (*types.Token, error) {
	parsedToken, err := methodutil.ParseJwtToken(accessToken, config.Jwt().AccessTokenSecret)
	if err != nil {
		return nil, errutil.ErrParseJwt
	}

	if _, ok := parsedToken.Claims.(jwt.Claims); !ok || !parsedToken.Valid {
		return nil, errutil.ErrInvalidAccessToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errutil.ErrInvalidAccessToken
	}

	return mapClaimsToToken(claims)
}
func (svc *TokenServiceImpl) ReadUserIDFromAccessTokenUUID(accessTokenUuid string) (int, error) {
	accessTokenCacheKey := config.Redis().MandatoryPrefix + config.Redis().AccessUuidPrefix + accessTokenUuid
	userID, err := svc.redisSvc.GetInt(accessTokenCacheKey)
	if err != nil {
		if errors.Is(err, errutil.ErrRecordNotFound) {
			return 0, errutil.ErrInvalidAccessToken
		}
		return 0, fmt.Errorf("error reading user ID from access token UUID: %v", err)
	}
	return userID, nil
}

func mapClaimsToToken(claims jwt.MapClaims) (*types.Token, error) {

	jsonData, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("error marshalling claims to JSON: %v", err)
	}
	var token types.Token
	if err := json.Unmarshal(jsonData, &token); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON to Token: %v", err)
	}
	return &token, nil
}
func (svc *TokenServiceImpl) DeleteTokenUUID(token *types.Token) error {
	err := svc.redisSvc.Delete(methodutil.AccessUuidCacheKey(token.AccessUuid), methodutil.RefreshUuidCacheKey(token.RefreshUuid))

	if err != nil {
		return err
	}

	return nil
}
