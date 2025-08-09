package cmd

import (
	"ems/config"
	"ems/conn"
	"ems/controllers"
	"ems/middlewares"
	asynq_repo "ems/repositories/asynq"
	db_repo "ems/repositories/db"
	"ems/repositories/mail"
	"ems/routes"
	"ems/server"
	"ems/services"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: Serve,
}

func Serve(cmd *cobra.Command, args []string) {
	config.LoadConfig()
	// clients
	dbClient := conn.Db()
	redisClient := conn.Redis()
	emailClient := conn.EmailClient()
	asynqClient := conn.Asynq()
	asynqInspector := conn.AsynqInspector()
	//worker
	workerPool := conn.WorkerPool()
	// repositories
	dbRepo := db_repo.NewRepository(dbClient)
	mailRepo := mail.NewRepository(emailClient, config.Email())
	asynqRepo := asynq_repo.NewRepository(config.Asynq(), asynqClient, asynqInspector)
	redisSvc := services.NewRedisService(redisClient)
	userSvc := services.NewUserServiceImpl(dbRepo, redisSvc)
	tokenSvc := services.NewTokenServiceImpl(redisSvc)
	authSvc := services.NewAuthServiceImpl(userSvc, tokenSvc)
	mailSvc := services.NewMailService(dbRepo, mailRepo, workerPool)
	asynqSvc := services.NewAsynqService(config.Asynq(), asynqRepo, dbRepo, dbRepo)
	campaignSvc := services.NewCampaignServiceImpl(dbRepo, mailSvc, asynqSvc)

	userCtrl := controllers.NewUserController(userSvc)
	authCtrl := controllers.NewAuthController(authSvc)
	campaignCtrl := controllers.NewCampaignController(campaignSvc)
	// middlewares
	authMiddleware := middlewares.NewAuthMiddleware(authSvc, userSvc)

	// Server
	echoServer := echo.New()
	srv := server.New(echoServer)

	routes := routes.New(echoServer, userCtrl, authCtrl, campaignCtrl, authMiddleware)

	routes.Init()
	srv.Start(workerPool)

}
