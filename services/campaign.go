package services

import (
	"ems/domain"
	"ems/models"
	"ems/types"
	"ems/utils/errutil"
	"ems/utils/msgutil"
	"errors"
	"fmt"
	"time"
)

type CampaignServiceImpl struct {
	repo     domain.CampaignRepository
	mailSvc  domain.MailService
	asynqSvc domain.AsynqService
}

func NewCampaignServiceImpl(campaignRepo domain.CampaignRepository, mailSvc domain.MailService, asynqSvc domain.AsynqService) *CampaignServiceImpl {
	return &CampaignServiceImpl{
		repo:     campaignRepo,
		mailSvc:  mailSvc,
		asynqSvc: asynqSvc,
	}
}
func (svc *CampaignServiceImpl) CreateCampaign(req *types.CampaignCreateRequest) (*types.CampaignCreateResponse, error) {
	// Check if a campaign with the same title already exists
	isExist, err := svc.IsCampaignExist(req.Title)
	if err != nil {
		return nil, err
	}
	if isExist {
		return nil, errutil.ErrAlreadyExists
	}

	// Convert request to campaign model
	campaign := req.ToCampaignModel()
	// Create campaign using repository
	createdCampaign, err := svc.repo.CreateCampaign(campaign)
	if err != nil {
		return nil, err
	}

	// Return successful response
	return &types.CampaignCreateResponse{
		Message:  msgutil.CampaignCreatedSuccessfully(),
		Campaign: createdCampaign,
	}, nil
}

// IsCampaignExist checks whether a campaign with the given title exists
func (svc *CampaignServiceImpl) IsCampaignExist(title string) (bool, error) {
	return svc.repo.ReadCampaignByTitle(title)
}

func (svc *CampaignServiceImpl) GetCampaignByID(campaignID int) (*types.CampaignCommonResponse, error) {
	campaign, err := svc.repo.ReadCampaignById(campaignID)
	if err != nil {
		return nil, err
	}
	return &types.CampaignCommonResponse{
		Message:  "Campaign Fetched successfully",
		Campaign: campaign,
	}, nil
}

func (svc *CampaignServiceImpl) UpdateCampaign(campaign *types.CampaignUpdateRequest, updatedBy int) (*types.CampaignUpdateResponse, error) {
	// Read campaign with ID and status "Draft"
	existingCampaign, err := svc.repo.ReadCampaignByIdAndStatus(campaign.ID, "Draft")
	if err != nil {
		return nil, err
	}
	if existingCampaign == nil {
		return nil, errutil.ErrRecordNotFound
	}
	req := campaign.ToCampaignModel()
	req.ID = existingCampaign.ID
	req.UpdatedBy = updatedBy
	updatedCampaign, err := svc.repo.UpdateCampaign(req)
	if err != nil {
		return nil, err
	}
	return &types.CampaignUpdateResponse{
		Message:  "Campaign Updated Successfully",
		Campaign: updatedCampaign,
	}, nil
}

func (svc *CampaignServiceImpl) DeleteCampaign(campaignID int) (*types.CampaignDeleteResponse, error) {
	err := svc.repo.DeleteCampaign(campaignID)
	if err != nil {
		return nil, err
	}
	return &types.CampaignDeleteResponse{
		Message: "Campaign Deleted Successfully",
	}, nil
}

func (svc *CampaignServiceImpl) ListCampaigns() (*types.CampaignCommonResponseList, error) {
	campaigns, err := svc.repo.ListCampaigns()
	if errors.Is(err, errutil.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &types.CampaignCommonResponseList{
		Message:  "Campaign List Fetched Successfully",
		Campaign: campaigns,
	}, nil
}

func (svc *CampaignServiceImpl) ApproveRejectCampaign(campaignID int, updatedBy int) (*types.CampaignApproveRejectResponse, error) {
	existingCampaign, err := svc.repo.ReadCampaignByIdAndStatus(campaignID, "Draft")
	if err != nil {
		return nil, err
	}
	if existingCampaign == nil {
		return nil, errutil.ErrRecordNotFound
	}
	err = svc.repo.ApproveRejectCampaign(campaignID, updatedBy)
	if err != nil {
		return nil, err
	}
	// Step 2: Audit Log
	err = svc.logApprovalActivity(existingCampaign, updatedBy)
	if err != nil {
		// Log error but don't fail approval
		fmt.Printf("Audit log failed for campaign Title %s: %v\n", existingCampaign.Title, err)
	}
	// Step 3: Send email notification
	go func() {
		err := svc.sendApprovalNotification(existingCampaign)
		if err != nil {
			fmt.Printf("Email notification failed for campaign Title %s: %v\n", existingCampaign.Title, err)
		}
	}()
	return &types.CampaignApproveRejectResponse{
		Message: "Event Successful",
	}, nil
}

func (svc *CampaignServiceImpl) logApprovalActivity(campaign *models.Campaign, updatedBy int) error {
	// You can store in a separate audit table or file
	logEntry := fmt.Sprintf("Campaign Title %d approved by %s at %s", campaign.Title, updatedBy, time.Now().Format(time.RFC3339))
	fmt.Println("[AUDIT]", logEntry)

	// TODO: Save to DB table if required
	return nil
}

func (svc *CampaignServiceImpl) sendApprovalNotification(campaign *models.Campaign) error {
	// Fetch campaign details (optional)
	var roleIds = []int{3}
	//err := svc.mailSvc.SendCampaignEmail(roleIds, campaign)
	err := svc.asynqSvc.AsynqTaskSendEmail(roleIds, campaign)
	return err
}
