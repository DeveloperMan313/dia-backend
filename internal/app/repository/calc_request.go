package repository

import (
	"dia-backend/internal/app/ds"
	"errors"

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

func (r *CalcRequestRepository) GetCalcRequestIDEntryCntByUserID(userID uint64) (int, int, error) {
	var calcRequest ds.CalcRequest
	err := r.db.
		Model(&ds.CalcRequest{}).
		Where("status = 1 AND user_id = ?", userID).
		Take(&calcRequest).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	var count int64
	err = r.db.
		Model(&ds.CalcRequest{}).
		Where("id = ?", calcRequest.ID).
		Joins("CalcRequestToLamp").
		Count(&count).Error
	if err != nil {
		return 0, 0, err
	}

	return int(calcRequest.ID), int(count), nil
}

func (r *CalcRequestRepository) GetCalcRequestByID(id uint64, userID uint64) (*ds.CalcRequest, error) {
	var calcRequest ds.CalcRequest
	err := r.db.
		Preload("CalcRequestToLamp").
		Preload("CalcRequestToLamp.Lamp").
		Where("status = 1 AND user_id = ?", userID).
		First(&calcRequest, id).Error
	if err != nil {
		return nil, err
	}

	return &calcRequest, nil
}

func (r *CalcRequestRepository) AddLampToCalcRequest(lampID uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var lamp ds.Lamp
		err := r.db.First(&lamp, lampID).Error
		if err != nil {
			return err
		}

		var calcRequest ds.CalcRequest
		err = r.db.
			Where("status = 1 AND user_id = ?", userID).
			Take(&calcRequest).Error
		notFound := errors.Is(err, gorm.ErrRecordNotFound)
		if err != nil && !notFound {
			return err
		}

		if notFound {
			calcRequest = ds.CalcRequest{
				User: ds.User{ID: userID},
			}
			err := r.db.Create(&calcRequest).Error
			if err != nil {
				return err
			}
		}

		calcRequestToLamp := ds.CalcRequestToLamp{
			RequestID: calcRequest.ID,
			LampID:    lampID,
		}
		r.db.Create(&calcRequestToLamp)

		return nil
	})
}

func (r *CalcRequestRepository) DeleteCalcRequest(requestID uint64, userID uint64) error {
	query := "UPDATE calc_requests SET status=2 WHERE status = 1 AND id = $1 AND user_id = $2;"
	err := r.db.Raw(query, requestID, userID).Row().Err()
	if err != nil {
		return err
	}

	return nil
}
