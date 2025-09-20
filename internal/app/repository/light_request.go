package repository

import (
	"dia-backend/internal/app/ds"
	"errors"

	"gorm.io/gorm"
)

type LightRequestRepository struct {
	db *gorm.DB
}

func NewLightRequestRepository(db *gorm.DB) *LightRequestRepository {
	return &LightRequestRepository{
		db: db,
	}
}

func (r *LightRequestRepository) GetLightRequestIDEntryCntByUserID(userID uint64) (int, int, error) {
	var lightRequest ds.LightRequest
	err := r.db.
		Model(&ds.LightRequest{}).
		Where("status = 1 AND user_id = ?", userID).
		Take(&lightRequest).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	var count int64
	err = r.db.
		Model(&ds.LightRequest{}).
		Where("id = ?", lightRequest.ID).
		Joins("LightRequestToLamp").
		Count(&count).Error
	if err != nil {
		return 0, 0, err
	}

	return int(lightRequest.ID), int(count), nil
}

func (r *LightRequestRepository) GetLightRequestByID(id uint64, userID uint64) (*ds.LightRequest, error) {
	var lightRequest ds.LightRequest
	err := r.db.
		Preload("LightRequestToLamp").
		Preload("LightRequestToLamp.Lamp").
		Where("status = 1 AND user_id = ?", userID).
		First(&lightRequest, id).Error
	if err != nil {
		return nil, err
	}

	return &lightRequest, nil
}

func (r *LightRequestRepository) AddLampToLightRequest(lampID uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var lamp ds.Lamp
		err := r.db.First(&lamp, lampID).Error
		if err != nil {
			return err
		}

		var lightRequest ds.LightRequest
		err = r.db.
			Where("status = 1 AND user_id = ?", userID).
			Take(&lightRequest).Error
		notFound := errors.Is(err, gorm.ErrRecordNotFound)
		if err != nil && !notFound {
			return err
		}

		if notFound {
			lightRequest = ds.LightRequest{
				User: ds.User{ID: userID},
			}
			err := r.db.Create(&lightRequest).Error
			if err != nil {
				return err
			}
		}

		lightRequestToLamp := ds.LightRequestToLamp{
			RequestID: lightRequest.ID,
			LampID:    lampID,
		}
		r.db.Create(&lightRequestToLamp)

		return nil
	})
}

func (r *LightRequestRepository) DeleteLightRequest(requestID uint64, userID uint64) error {
	query := "UPDATE light_requests SET status=2 WHERE status = 1 AND id = $1 AND user_id = $2;"
	err := r.db.Raw(query, requestID, userID).Row().Err()
	if err != nil {
		return err
	}

	return nil
}
