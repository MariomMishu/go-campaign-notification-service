package controllers

import (
	"ems/consts"
	"ems/domain"
	"ems/middlewares"
	"ems/types"
	"ems/utils/errutil"
	"ems/utils/msgutil"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userSvc domain.UserService
}

func NewUserController(userSvc domain.UserService) *UserController {
	return &UserController{
		userSvc: userSvc,
	}
}

func (ctrl *UserController) SignUp(c echo.Context) error {
	var req types.CreateUserReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// if err := c.Validate(&req); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	// }
	req.RoleId = consts.RoleIdCustomer
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	if err := ctrl.userSvc.CreateUser(&req); err != nil {
		if errors.Is(err, errutil.ErrUserIsAlreadyExists) {
			return c.JSON(http.StatusConflict, msgutil.UserAlreadyExists())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusCreated, msgutil.UserCreatedSuccessfully())
}
func (ctrl *UserController) GetProfile(c echo.Context) error {
	currentUser, err := middlewares.CurrentUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
	}

	user, err := ctrl.userSvc.ReadUser(currentUser.ID)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, &types.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      consts.RoleMap[user.RoleID],
		RoleID:    user.RoleID,
	})
}

func (ctrl *UserController) CreateUser(c echo.Context) error {
	var req types.CreateUserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}
	if err := ctrl.userSvc.CreateUser(&req); err != nil {
		if errors.Is(err, errutil.ErrUserIsAlreadyExists) {
			return c.JSON(http.StatusConflict, msgutil.UserAlreadyExists())
		}
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}
	return c.JSON(http.StatusCreated, msgutil.UserCreatedSuccessfully())
}
