package cmd

import (
	"ems/config"
	"ems/conn"
	"ems/controllers"
	asynq_repo "ems/repositories/asynq"
	db_repo "ems/repositories/db"
	mail_Repo "ems/repositories/mail"
	"ems/services"
	"ems/types"
	"ems/worker"
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use: "worker",
	Run: runWorker,
}

func runWorker(cmd *cobra.Command, args []string) {
	//client
	dbClient := conn.Db()
	emailClient := conn.EmailClient()
	//worker
	workerPool := conn.WorkerPool()
	//repositories
	dbRepo := db_repo.NewRepository(dbClient)
	mailRepo := mail_Repo.NewRepository(emailClient, config.Email())
	asynqRepo := asynq_repo.NewRepository(config.Asynq())

	// services
	mailSvc := services.NewMailService(dbRepo, mailRepo, workerPool)
	asynqSvc := services.NewAsynqService(config.Asynq(), asynqRepo, dbRepo, dbRepo, mailSvc)

	// controllers
	asynqCtrl := controllers.NewAsynqController(mailSvc, asynqSvc)

	//services
	services.NewMailService(dbRepo, mailRepo, workerPool)
	mux := asynq.NewServeMux()
	mux.HandleFunc(types.AsynqTaskTypeSendEmail.String(), asynqCtrl.ProcessSendEmailTask)
	worker.StartAsynqWorker(mux)
}
