package conn

import (
	"ems/config"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

var emailClient *http.Client

func ConnectEmail() {
	conf := config.Email()
	if conf == nil {
		log.Error("Email config is nil!")
		return
	}
	log.Info("connecting to Email")
	timeout := conf.Timeout * time.Second
	emailClient = newHTTPClient(timeout, 50)
}

func EmailClient() *http.Client {
	return emailClient
}
