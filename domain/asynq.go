package domain

import (
	"ems/models"
	"ems/types"
	"github.com/hibiken/asynq"
)

type (
	AsynqRepository interface {
		CreateTask(campaign types.AsynqTaskType, payload interface{}) (*asynq.Task, error)
		EnqueueTask(task *asynq.Task, customOpts *types.AsynqOption) (string, error)
		DequeueTask(taskID string) error
	}

	AsynqService interface {
		AsynqTaskSendEmail(roleIds []int, campaign *models.Campaign) error
	}
)
