package repository

import (
	"dia-backend/internal/app/ds"

	"gorm.io/gorm"
)

type CalcRequestRepository struct {
	db *gorm.DB
}

func NewCalcRequestRepository(db *gorm.DB) *CalcRequestRepository {
	return &CalcRequestRepository{
		db: db,
	}
}

func (r *CalcRequestRepository) GetCalcRequestEntryCntByID(id uint64) (int64, error) {
	var count int64
	err := r.db.Model(&ds.CalcRequest{}).Where("id = ?", id).Joins("CalcRequestToLamp").Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *CalcRequestRepository) GetCalcRequestByID(id uint64, lampRepo *LampRepository) (*ds.CalcRequest, error) {
	var calcRequest ds.CalcRequest
	err := r.db.
		Preload("CalcRequestToLamp").
		Preload("CalcRequestToLamp.Lamp").
		Where("id = ?", id).
		Take(&calcRequest).Error
	if err != nil {
		return nil, err
	}

	return &calcRequest, nil
}
