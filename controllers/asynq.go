package controllers

import (
	"context"
	"ems/domain"
	"ems/types"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/labstack/gommon/log"
)

type AsynqController struct {
	mailSvc  domain.MailService
	asynqSvc domain.AsynqService
}

func NewAsynqController(mailSvc domain.MailService, asynqSvc domain.AsynqService) *AsynqController {
	return &AsynqController{
		mailSvc:  mailSvc,
		asynqSvc: asynqSvc,
	}
}
func (ac *AsynqController) ProcessSendEmailTask(ctx context.Context, t *asynq.Task) (err error) {
	log.Info(fmt.Sprintf("Received task event [%s] with ID [%s]", t.Type(), t.ResultWriter().TaskID()))
	var payload types.EmailPayload

	if err = json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Error(err)
		return
	}

	if err = ac.mailSvc.SendEmail(payload); err != nil {
		log.Error(fmt.Sprintf("err: [%v] occurred while sending email to: %s", err, payload.MailTo))
		return err
	}
	t.ResultWriter().Write([]byte(fmt.Sprintf("Email sent successfully to %s", payload.MailTo)))
	return
}
