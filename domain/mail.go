package domain

import (
	"ems/models"
	"ems/types"
)

type (
	MailService interface {
		SendEmail(requestData types.EmailPayload) error
		SendCampaignEmail(roleIds []int, campaign *models.Campaign) error
	}

	MailRepository interface {
		SendEmail(requestData *types.EmailPayload) error
	}
)
