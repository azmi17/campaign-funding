package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error) {
	transaction := []Transaction{}

	err := r.db.Preload("User").Where("campaign_id = ?", input.ID).Order("id desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, err
}
