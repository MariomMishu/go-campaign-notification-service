package db

import (
	"ems/models"
	"ems/utils/errutil"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func (repo *Repository) CreateCampaign(campaign *models.Campaign) (*models.Campaign, error) {
	campaign.Status = "Draft"
	qry := repo.db.Create(campaign)
	if qry.Error != nil {
		return nil, qry.Error
	}
	return campaign, nil
}

func (repo *Repository) ReadCampaignByTitle(title string) (bool, error) {
	var count int64
	if err := repo.db.Model(&models.Campaign{}).Where("title = ?", title).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *Repository) ReadCampaignById(id int) (*models.Campaign, error) {
	var campaign models.Campaign
	if err := repo.db.Model(&models.Campaign{}).Where("id = ?", id).First(&campaign).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (repo *Repository) ReadCampaignByIdAndStatus(id int, status string) (*models.Campaign, error) {
	var campaign models.Campaign
	if err := repo.db.
		Model(&models.Campaign{}).
		Where("id = ? AND status = ?", id, status).
		First(&campaign).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (repo *Repository) UpdateCampaign(campaign *models.Campaign) (*models.Campaign, error) {
	qry := repo.db.Where("id = ?", campaign.ID).Updates(&campaign)
	if errors.Is(qry.Error, gorm.ErrRecordNotFound) {
		return nil, errutil.ErrRecordNotFound
	}
	if qry.Error != nil {
		return nil, qry.Error
	}
	return campaign, nil
}

func (repo *Repository) DeleteCampaign(id int) error {
	qry := repo.db.Where("id = ?", id).Delete(&models.Campaign{})
	if qry.Error != nil {
		return qry.Error
	}
	if qry.RowsAffected == 0 {
		return errutil.ErrRecordNotFound
	}
	return nil
}
func (repo *Repository) ListCampaigns() ([]*models.Campaign, error) {
	var campaigns []*models.Campaign
	if err := repo.db.Model(&models.Campaign{}).Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (repo *Repository) ApproveRejectCampaign(id int, updatedBy int) error {
	// Step 1: Update campaign status

	qry := repo.db.Model(&models.Campaign{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "Approved",
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	})
	if qry.Error != nil {
		return qry.Error
	}
	if qry.RowsAffected == 0 {
		return errutil.ErrRecordNotFound
	}
	fmt.Printf("Campaign ID %d approved by %s at %s\n", id, 1, time.Now().Format(time.RFC3339))

	return nil
}
