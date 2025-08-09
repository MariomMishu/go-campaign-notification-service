package types

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	CreateUserReq struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RoleId    int    `json:"role_id"`
	}
	UserInfo struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RoleID    int    `json:"role_id"`
		Role      string `json:"role,omitempty" gorm:"-"`
	}

	CurrentUser struct {
		ID          int    `json:"id"`
		Email       string `json:"email"`
		RoleID      int    `json:"role_id"`
		Role        string `json:"role"`
		AccessUuid  string `json:"access_uuid"`
		RefreshUuid string `json:"refresh_uuid"`
	}
)

func (crq *CreateUserReq) Validate() error {
	return v.ValidateStruct(crq,
		v.Field(&crq.Email, v.Required, is.Email),
		v.Field(&crq.Password, v.Required),
		v.Field(&crq.FirstName, v.Required, v.Length(0, 50)),
		v.Field(&crq.LastName, v.Required, v.Length(0, 50)),
		v.Field(&crq.RoleId, v.Required),
	)
}
