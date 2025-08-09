package routes

import (
	"ems/consts"
	"ems/controllers"
	m "ems/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Route struct {
	echo           *echo.Echo
	userCtrl       *controllers.UserController
	authCtrl       *controllers.AuthController
	campaignCtrl   *controllers.CampaignController
	authMiddleware *m.AuthMiddleware
}

func New(e *echo.Echo, userCtrl *controllers.UserController, authCtrl *controllers.AuthController, campaignCtrl *controllers.CampaignController, authMiddleware *m.AuthMiddleware) *Route {
	return &Route{
		echo:           e,
		userCtrl:       userCtrl,
		authCtrl:       authCtrl,
		campaignCtrl:   campaignCtrl,
		authMiddleware: authMiddleware,
	}
}

func (r *Route) Init() {
	e := r.echo
	m.Init(e) //Middleware Initialization
	// APM routes
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	g := e.Group("/v1")
	user := g.Group("/user")
	campaign := g.Group("/campaign")

	user.POST("/create", r.userCtrl.SignUp)
	user.GET("/profile", r.userCtrl.GetProfile, r.authMiddleware.Authenticate(""))
	user.POST("", r.userCtrl.CreateUser, r.authMiddleware.Authenticate(consts.PermissionUserCreate))

	auth := g.Group("/auth")
	auth.POST("/login", r.authCtrl.Login)
	auth.POST("/logout", r.authCtrl.Logout, r.authMiddleware.Authenticate(""))

	campaign.POST("", r.campaignCtrl.CreateCampaign, r.authMiddleware.Authenticate(consts.PermissionCampaignCreate))
	campaign.GET("", r.campaignCtrl.GetCampaignList, r.authMiddleware.Authenticate(consts.PermissionCampaignList))
	campaign.GET("/:id", r.campaignCtrl.GetCampaignById, r.authMiddleware.Authenticate(consts.PermissionCampaignFetch))
	campaign.PUT("/:id", r.campaignCtrl.UpdateCampaign, r.authMiddleware.Authenticate(consts.PermissionCampaignUpdate))
	campaign.DELETE("/:id", r.campaignCtrl.DeleteCampaign, r.authMiddleware.Authenticate(consts.PermissionCampaignDelete))
	campaign.PUT("/action/:id", r.campaignCtrl.ApproveRejectCampaign, r.authMiddleware.Authenticate(consts.PermissionCampaignApproveReject))

}
