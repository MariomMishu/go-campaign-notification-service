package cmd

import (
	"ems/config"
	"ems/conn"
	db_repo "ems/repositories/db"
	mail_Repo "ems/repositories/mail"
	"ems/services"
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

	//services
	services.NewMailService(dbRepo, mailRepo, workerPool)
	mux := asynq.NewServeMux()
	worker.StartAsynqWorker(mux)
}
