package cmd

import (
	"ems/config"
	"ems/conn"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	RootCmd = &cobra.Command{
		Use: "app",
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
	RootCmd.AddCommand(workerCmd)
}

// Execute executes the root command
func Execute() {
	// load config
	config.LoadConfig()

	// Initialize logger
	initLogger()

	conn.ConnectDb()
	conn.ConnectRedis()
	conn.ConnectEmail()
	conn.ConnectWorker()
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initLogger() {
	fmt.Println("Initializing logger...")
}
