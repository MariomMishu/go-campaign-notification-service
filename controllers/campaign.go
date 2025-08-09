package controllers

import (
	"ems/domain"
	"ems/middlewares"
	"ems/types"
	"ems/utils/errutil"
	"ems/utils/msgutil"
	"errors"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CampaignController struct {
	campaignSvc domain.CampaignService
}

func NewCampaignController(campaignSvc domain.CampaignService) *CampaignController {
	return &CampaignController{
		campaignSvc: campaignSvc,
	}
}

func (ctrl *CampaignController) CreateCampaign(c echo.Context) error {
	var req types.CampaignCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}
	if err := v.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.SomethingWentWrongMsg())
	}
	currentUser, err := middlewares.CurrentUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
	}
	req.CreatedBy = currentUser.ID
	resp, err := ctrl.campaignSvc.CreateCampaign(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

func (ctrl *CampaignController) DeleteCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id")) // Extracts the id parameter from the URL and Converts the string "12" to an integer 12 using strconv.Atoi.
	if err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg()) // If the ID is not a valid number, return a 400 Bad Request with a custom message.
	}
	if err := v.Validate(id, v.Required); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}
	resp, err := ctrl.campaignSvc.DeleteCampaign(id)
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, msgutil.CampaignNotFound) //If the record doesn’t exist in the database, return 404 Not Found.
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg()) //If any other error occurs (like DB connection issues), return 500 Internal Server Error.
	}
	return c.JSON(http.StatusOK, resp)
}

func (ctrl *CampaignController) UpdateCampaign(c echo.Context) error {
	var req types.CampaignUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}
	if err := v.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}
	user, err := middlewares.CurrentUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
	}
	resp, err := ctrl.campaignSvc.UpdateCampaign(&req, user.ID)
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, msgutil.CampaignNotFound)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}
	return c.JSON(http.StatusOK, resp)
}

func (ctrl *CampaignController) GetCampaignList(c echo.Context) error {
	campaigns, err := ctrl.campaignSvc.ListCampaigns()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}
	return c.JSON(http.StatusOK, campaigns)
}

func (ctrl *CampaignController) GetCampaignById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id")) // Extracts the id parameter from the URL and Converts the string "12" to an integer 12 using strconv.Atoi.
	if err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg()) // If the ID is not a valid number, return a 400 Bad Request with a custom message.
	}
	if err := v.Validate(id, v.Required); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}
	resp, err := ctrl.campaignSvc.GetCampaignByID(id)
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, msgutil.CampaignNotFound) //If the record doesn’t exist in the database, return 404 Not Found.
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg()) //If any other error occurs (like DB connection issues), return 500 Internal Server Error.
	}
	return c.JSON(http.StatusOK, resp)
}

func (ctrl *CampaignController) ApproveRejectCampaign(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id")) // Extracts the id parameter from the URL and Converts the string "12" to an integer 12 using strconv.Atoi.
	if err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg()) // If the ID is not a valid number, return a 400 Bad Request with a custom message.
	}
	user, err := middlewares.CurrentUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, msgutil.UserUnauthorized())
	}
	resp, err := ctrl.campaignSvc.ApproveRejectCampaign(id, user.ID)
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, msgutil.CampaignNotFound)
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}
	return c.JSON(http.StatusOK, resp)
}
