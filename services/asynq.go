package services

import (
	"ems/config"
	"ems/domain"
	"ems/models"
	"ems/types"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/labstack/gommon/log"
)

type AsynqService struct {
	config       *config.AsynqConfig
	asynqRepo    domain.AsynqRepository
	userRepo     domain.UserRepository
	campaignRepo domain.CampaignRepository
	mailSvc      domain.MailService
}

func NewAsynqService(
	config *config.AsynqConfig,
	asynqRepo domain.AsynqRepository,
	userRepo domain.UserRepository,
	campaignRepo domain.CampaignRepository,
	mailSvc domain.MailService,
) *AsynqService {
	return &AsynqService{
		config:       config,
		asynqRepo:    asynqRepo,
		userRepo:     userRepo,
		campaignRepo: campaignRepo,
		mailSvc:      mailSvc,
	}
}

func (svc *AsynqService) AsynqTaskSendEmail(roleIds []int, campaign *models.Campaign) error {
	users, err := svc.userRepo.ReadUsers(roleIds)
	if err != nil {
		log.Printf("Failed to read users: %v", err)
		return err
	}

	for _, user := range users {
		task, err := svc.sendEmailCampaignApprovalTask(user, campaign)

		if err != nil {
			log.Error(fmt.Sprintf("err: [%v] occurred while creating email sending task for user: %v", err, user.Email))
			return err
		}

		taskID := fmt.Sprintf("%s_user:%d_campaign:%d", types.AsynqTaskTypeSendEmail, user.ID, campaign.ID)
		customOpts := &types.AsynqOption{
			Queue:        svc.config.Queue,
			TaskID:       taskID,
			DelaySeconds: svc.config.EmailSendTaskDelay,
			Retry:        svc.config.EmailSendTaskRetryCount,
		}
		_, err = svc.enqueueTask(task, customOpts)
		if err != nil {
			return err
		}
		log.Info(fmt.Sprintf("enqueued email send task for user [%s] successfully", user.Email))
	}

	return nil
}

func (svc *AsynqService) sendEmailCampaignApprovalTask(user *models.User, campaign *models.Campaign) (*asynq.Task, error) {
	emailPayload := types.EmailPayload{
		MailTo:  user.Email,
		Subject: "Test Email For Campaign",
		Body: map[string]interface{}{
			"Title":       campaign.Title,
			"Description": campaign.Description,
			"Remarks":     campaign.Remarks,
		},
	}
	svc.mailSvc.SendEmail(emailPayload)
	return svc.asynqRepo.CreateTask(types.AsynqTaskTypeSendEmail, emailPayload)
}
func (svc *AsynqService) enqueueTask(task *asynq.Task, customOpts *types.AsynqOption) (taskID string, err error) {
	err = svc.asynqRepo.DequeueTask(customOpts.TaskID) // Ensure no duplicate tasks
	if err != nil && !errors.Is(err, asynq.ErrTaskNotFound) {
		log.Error(fmt.Sprintf("error: [%v] occurred while dequeuing task with ID: %s", err, customOpts.TaskID))
	}

	taskID, err = svc.asynqRepo.EnqueueTask(task, customOpts)
	if errors.Is(err, asynq.ErrDuplicateTask) {
		log.Warn(fmt.Sprintf("skipped: duplicate task for taskID: [%s]", customOpts.TaskID))
		err = nil // No error for duplicate tasks, just skip
		return
	}
	if err != nil {
		log.Error(fmt.Sprintf("error: [%v] occurred while enqueuing task with ID: %s", err, customOpts.TaskID))
		return
	}

	log.Info(fmt.Sprintf("enqueued task [%s] successfully", taskID))
	return taskID, nil
}
