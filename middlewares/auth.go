package middlewares

import (
	"ems/domain"
	"ems/models"
	"ems/types"
	"fmt"
	"strings"

	"ems/utils/errutil"
	"ems/utils/msgutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

const ContextKeyCurrentUser = "user"

type AuthMiddleware struct {
	authSvc domain.AuthService
	userSvc domain.UserService
}

func NewAuthMiddleware(authSvc domain.AuthService, userSvc domain.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (m *AuthMiddleware) Authenticate2() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString, err := m.tokenFromHeader(c)

			if err != nil {

				return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
			}

			userInfo, token, err := m.authSvc.VerifyAccessToken(tokenString)
			if err != nil {
				fmt.Println("tokenString: http.StatusUnauthorized")
				return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
			}

			currentUser := types.CurrentUser{
				ID:          userInfo.ID,
				Email:       userInfo.Email,
				RoleID:      userInfo.RoleID,
				Role:        userInfo.Role,
				AccessUuid:  token.AccessUuid,
				RefreshUuid: token.RefreshUuid,
			}

			// Set user in context
			c.Set(ContextKeyCurrentUser, currentUser)
			// Set user ID and permissions in header
			c.Request().Header.Set("X-User-ID", fmt.Sprintf("%d", userInfo.ID))

			return next(c)
		}
	}
}

func (m *AuthMiddleware) Authenticate(permission string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString, err := m.tokenFromHeader(c)

			if err != nil {

				return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
			}

			userInfo, token, err := m.authSvc.VerifyAccessToken(tokenString)
			if err != nil {
				fmt.Println("tokenString: http.StatusUnauthorized")
				return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
			}

			currentUser := types.CurrentUser{
				ID:          userInfo.ID,
				Email:       userInfo.Email,
				RoleID:      userInfo.RoleID,
				Role:        userInfo.Role,
				AccessUuid:  token.AccessUuid,
				RefreshUuid: token.RefreshUuid,
			}

			// Set user in context
			c.Set(ContextKeyCurrentUser, currentUser)
			// Set user ID and permissions in header
			permissions, err := m.userSvc.ReadPermissionsByRole(currentUser.RoleID)
			//for _, permission := range permissions {
			//	fmt.Println(permission.Permission)
			//}
			if err != nil {

				return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
			}
			if permission != "" && !m.isPermissionAllowed(permission, permissions) {
				return c.JSON(http.StatusForbidden, msgutil.AccessForbiddenMsg())
			}
			c.Request().Header.Set("X-User-ID", fmt.Sprintf("%d", userInfo.ID))

			return next(c)
		}
	}
}
func (m *AuthMiddleware) tokenFromHeader(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", errutil.ErrInvalidAuthorizationToken
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == "" {
		return "", errutil.ErrInvalidAuthorizationToken
	}

	return tokenString, nil
}
func (m *AuthMiddleware) isPermissionAllowed(permission string, rolePermissions []*models.Permission) bool {
	for _, rolePermission := range rolePermissions {
		if rolePermission.Permission == permission {
			return true
		}
	}
	return false
}
