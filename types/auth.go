package types

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type (
	LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginResp struct {
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		User         *UserInfo `json:"user"`
	}

	// UserInfo represents user information returned in LoginResp
	// UserInfo struct {
	// 	ID    string `json:"id"`
	// 	Email string `json:"email"`
	// 	// Add other fields as needed
	// }
)

func (req *LoginReq) Validate() error {
	return v.ValidateStruct(req,
		v.Field(&req.Email, v.Required, is.Email),
		v.Field(&req.Password, v.Required),
	)
}
