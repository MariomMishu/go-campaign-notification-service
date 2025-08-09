package types

import (
	"ems/models"
	"time"
)

type (
	CampaignCreateRequest struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Remarks     string  `json:"remarks"`
		Status      string  `json:"status"`
		StartTime   *string `json:"start_time"` // ISO8601 format recommended
		EndTime     *string `json:"end_time"`   // ISO8601 format recommended
		CreatedBy   int     `json:"created_by"`
	}
	CampaignCreateResponse struct {
		Campaign *models.Campaign `json:"campaign"`
		Message  interface{}      `json:"message"`
	}
	CampaignUpdateRequest struct {
		ID int `param:"id"`
		CampaignCreateRequest
	}
	CampaignDeleteResponse struct {
		Message string `json:"message"`
	}
	CampaignApproveRejectResponse struct {
		Message string `json:"message"`
	}
	CampaignUpdateResponse struct {
		Message  string           `json:"message"`
		Campaign *models.Campaign `json:"campaign"`
	}
	CampaignCommonResponse struct {
		Campaign *models.Campaign `json:"campaign"`
		Message  string           `json:"message"`
	}
	CampaignCommonResponseList struct {
		Campaign []*models.Campaign `json:"campaigns"`
		Message  string             `json:"message"`
	}
)

// If you need the second struct, rename it, for example:

func (createReq *CampaignCreateRequest) ToCampaignModel() *models.Campaign {
	campaign := &models.Campaign{
		Title:       createReq.Title,
		Description: createReq.Description,
		Remarks:     createReq.Remarks,
		CreatedBy:   createReq.CreatedBy,
	}
	if createReq.StartTime != nil {
		campaign.StartTime, _ = parseTime(*createReq.StartTime, time.RFC3339)
	}
	if createReq.EndTime != nil {
		campaign.EndTime, _ = parseTime(*createReq.EndTime, time.RFC3339)
	}
	return campaign
}

func (updateReq *CampaignUpdateRequest) ToCampaignModel() *models.Campaign {
	campaign := &models.Campaign{
		Title:       updateReq.Title,
		Description: updateReq.Description,
		Remarks:     updateReq.Remarks,
		CreatedBy:   updateReq.CreatedBy,
	}
	if updateReq.StartTime != nil {
		campaign.StartTime, _ = parseTime(*updateReq.StartTime, time.RFC3339)
	}
	if updateReq.EndTime != nil {
		campaign.EndTime, _ = parseTime(*updateReq.EndTime, time.RFC3339)
	}
	return campaign
}

func parseTime(timeStr string, format string) (*time.Time, error) {
	t, err := time.Parse(format, timeStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
