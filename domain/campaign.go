package domain

import (
	"ems/models"
	"ems/types"
)

// Campaign represents an campaign entity.
type Campaign struct {
	ID          int
	Title       string
	Description string
	StartTime   string // Consider using time.Time in real applications
	EndTime     string
	Remarks     string
	Status      string
}

type (
	CampaignRepository interface {
		CreateCampaign(campaign *models.Campaign) (*models.Campaign, error)
		ReadCampaignById(id int) (*models.Campaign, error)
		ReadCampaignByIdAndStatus(id int, status string) (*models.Campaign, error)
		ReadCampaignByTitle(title string) (bool, error)
		UpdateCampaign(campaign *models.Campaign) (*models.Campaign, error)
		DeleteCampaign(id int) error
		ApproveRejectCampaign(id int, updatedBy int) error
		ListCampaigns() ([]*models.Campaign, error)
	}
	CampaignService interface {
		CreateCampaign(campaign *types.CampaignCreateRequest) (*types.CampaignCreateResponse, error)
		GetCampaignByID(id int) (*types.CampaignCommonResponse, error)
		UpdateCampaign(campaign *types.CampaignUpdateRequest, updatedBy int) (*types.CampaignUpdateResponse, error)
		DeleteCampaign(id int) (*types.CampaignDeleteResponse, error)
		ApproveRejectCampaign(id int, updatedBy int) (*types.CampaignApproveRejectResponse, error)
		ListCampaigns() (*types.CampaignCommonResponseList, error)
	}
)
