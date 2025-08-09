package mail

import (
	"ems/config"
	"ems/types"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

type Repository struct {
	client *http.Client
	config *config.EmailConfig
}

func NewRepository(client *http.Client, config *config.EmailConfig) *Repository {
	return &Repository{
		config: config,
		client: client,
	}
}

//	func (r *Repository) SendEmail(requestData *types.EmailPayload) error {
//		reqURL := config.Email().Url
//
//		reqByte, err := json.Marshal(requestData)
//		if err != nil {
//			log.Printf("Failed to marshal email payload: %v", err)
//			return err
//		}
//
//		req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(reqByte))
//		if err != nil {
//			log.Printf("Failed to create email request: %v", err)
//			return err
//		}
//		req.Header.Set("Content-Type", "application/json")
//
//		resp, err := r.client.Do(req)
//		if err != nil {
//			log.Printf("Failed to send email request: %v", err)
//			return err
//		}
//		defer resp.Body.Close()
//
//		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
//			log.Printf("Email service returned non-success status: %v", resp.Status)
//			return fmt.Errorf("email service returned non-success status: %v", resp.Status)
//		}
//
//		log.Infof("Email service responded with status [%v] for recipient: %v", resp.StatusCode, requestData.MailTo)
//		return nil
//	}
func (r *Repository) SendEmail(payload *types.EmailPayload) error {
	emailCfg := config.Email()
	mailTo := strings.Split(payload.MailTo, "@")
	mailFrom := emailCfg.Username
	auth := smtp.PlainAuth("", emailCfg.Username, emailCfg.Password, emailCfg.Host)

	// Build the email message
	msg := "Test Email"

	addr := emailCfg.Host + ":" + emailCfg.Port

	// Optional timeout handling
	done := make(chan error, 1)
	go func() {
		err := smtp.SendMail(addr, auth, mailFrom, mailTo, []byte(msg))
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
		return nil
	case <-time.After(emailCfg.Timeout):
		return fmt.Errorf("email send timeout after %v", emailCfg.Timeout)
	}
}
