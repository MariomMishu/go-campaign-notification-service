package methodutil

import (
	"ems/config"
	"ems/utils/errutil"
	"github.com/golang-jwt/jwt/v4"
)

func ParseJwtToken(tokenString, secret string) (*jwt.Token, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errutil.ErrInvalidJwtSigningMethod
		}
		return []byte(secret), nil
	}

	return jwt.Parse(tokenString, keyFunc)

}
func AccessUuidCacheKey(accessUuid string) string {
	return config.Redis().MandatoryPrefix + config.Redis().AccessUuidPrefix + accessUuid
}

func RefreshUuidCacheKey(refreshUuid string) string {
	return config.Redis().MandatoryPrefix + config.Redis().RefreshUuidPrefix + refreshUuid
}
