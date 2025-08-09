package server

import (
	"context"
	"ems/config"
	"ems/worker"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	echo *echo.Echo
}

func New(echo *echo.Echo) *Server {
	return &Server{echo: echo}
}
func (s *Server) Start(workerPool *worker.Pool) {
	e := s.echo
	go func() {
		if err := e.Start(":" + config.App().Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Server shutdown with error: %v", err)
	}
	log.Info("Server exited gracefully")
	workerPool.StopWithContext(ctx)
}
