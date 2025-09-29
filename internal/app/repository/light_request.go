package repository

import (
	"dia-backend/internal/app/ds"
	"errors"
	"time"

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

func (r *LightRequestRepository) GetDraftRequestInfo(userID uint64) (uint64, int, error) {
	var lightRequest ds.LightRequest
	err := r.db.
		Where("status = 1 AND user_id = ?", userID).
		First(&lightRequest).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	var count int64
	err = r.db.
		Model(&ds.LightRequestToLamp{}).
		Where("request_id = ?", lightRequest.ID).
		Count(&count).Error

	if err != nil {
		return 0, 0, err
	}

	return lightRequest.ID, int(count), nil
}

func (r *LightRequestRepository) GetLightRequests(statusFilter uint8, dateFrom, dateTo *time.Time) ([]ds.LightRequest, error) {
	var lightRequests []ds.LightRequest

	query := r.db.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Preload("Moderator", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Where("status != 1 AND status != 2")

	if statusFilter != 0 {
		query = query.Where("status = ?", statusFilter)
	}

	if dateFrom != nil {
		query = query.Where("formed_at >= ?", dateFrom)
	}
	if dateTo != nil {
		query = query.Where("formed_at <= ?", dateTo)
	}

	err := query.Find(&lightRequests).Error
	if err != nil {
		return nil, err
	}

	return lightRequests, nil
}

func (r *LightRequestRepository) GetLightRequestByID(id uint64, userID uint64) (*ds.LightRequest, error) {
	var lightRequest ds.LightRequest
	err := r.db.
		Preload("LightRequestToLamp").
		Preload("LightRequestToLamp.Lamp").
		Where("status != 2 AND user_id = ?", userID).
		First(&lightRequest, id).Error

	if err != nil {
		return nil, err
	}

	return &lightRequest, nil
}

func (r *LightRequestRepository) UpdateLightRequest(id uint64, maxTotalPowerW *float64) error {
	updates := make(map[string]interface{})

	if maxTotalPowerW != nil {
		updates["max_total_power_w"] = *maxTotalPowerW
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.
		Model(&ds.LightRequest{}).
		Where("id = ? AND status != 2", id).
		Updates(updates).Error
}

func (r *LightRequestRepository) FormRequest(id uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var lightRequest ds.LightRequest
		err := tx.
			Preload("LightRequestToLamp").
			Where("id = ? AND user_id = ? AND status = 1", id, userID).
			First(&lightRequest).Error

		if err != nil {
			return err
		}

		if lightRequest.MaxTotalPowerW <= 0 {
			return errors.New("max total power should be positive")
		}
		if len(lightRequest.LightRequestToLamp) == 0 {
			return errors.New("at least one lamp is required")
		}

		return tx.Model(&lightRequest).Updates(map[string]interface{}{
			"status":    3,
			"formed_at": time.Now(),
		}).Error
	})
}

func (r *LightRequestRepository) ResolveOrRejectRequest(id uint64, moderatorID uint64, status uint8) error {
	if status != 4 && status != 5 {
		return errors.New("invalid status for moderator action")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		var lightRequest ds.LightRequest
		err := tx.
			Where("id = ? AND status = 3", id).
			First(&lightRequest).Error

		if err != nil {
			return err
		}

		updates := map[string]interface{}{
			"status":       status,
			"moderator_id": moderatorID,
			"closed_at":    time.Now(),
		}

		return tx.Model(&lightRequest).Updates(updates).Error
	})
}

func (r *LightRequestRepository) DeleteRequest(id uint64, userID uint64) error {
	return r.db.
		Model(&ds.LightRequest{}).
		Where("id = ? AND user_id = ? AND status = 1", id, userID).
		Update("status", 2).Error
}

func (r *LightRequestRepository) RemoveLampFromRequest(requestID uint64, lampID uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var lightRequest ds.LightRequest
		err := tx.
			Where("id = ? AND user_id = ? AND status = 1", requestID, userID).
			First(&lightRequest).Error

		if err != nil {
			return err
		}

		return tx.
			Where("request_id = ? AND lamp_id = ?", requestID, lampID).
			Delete(&ds.LightRequestToLamp{}).Error
	})
}

func (r *LightRequestRepository) UpdateRequestToLamp(requestID uint64, lampID uint64, userID uint64, areaM2 *float64, number *uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var lightRequest ds.LightRequest
		err := tx.
			Where("id = ? AND user_id = ? AND status = 1", requestID, userID).
			First(&lightRequest).Error

		if err != nil {
			return err
		}

		updates := make(map[string]interface{})
		if areaM2 != nil {
			updates["area_m2"] = *areaM2
		}
		if number != nil {
			updates["number"] = *number
		}

		if len(updates) == 0 {
			return nil
		}

		return tx.
			Model(&ds.LightRequestToLamp{}).
			Where("request_id = ? AND lamp_id = ?", requestID, lampID).
			Updates(updates).Error
	})
}

func (r *LightRequestRepository) CalculateTotalPower(requestID uint64) float64 {
	var result struct {
		TotalPower float64
	}

	r.db.Model(&ds.LightRequestToLamp{}).
		Select("SUM(lamps.power_w * light_request_to_lamps.number) as total_power").
		Joins("JOIN lamps ON lamps.id = light_request_to_lamps.lamp_id").
		Where("light_request_to_lamps.request_id = ?", requestID).
		Scan(&result)

	return result.TotalPower
}

func (r *LightRequestRepository) AddLampToLightRequest(lampID uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var lamp ds.Lamp
		err := tx.First(&lamp, lampID).Error
		if err != nil {
			return err
		}

		var lightRequest ds.LightRequest
		err = tx.
			Where("status = 1 AND user_id = ?", userID).
			Take(&lightRequest).Error

		notFound := errors.Is(err, gorm.ErrRecordNotFound)
		if err != nil && !notFound {
			return err
		}

		if notFound {
			lightRequest = ds.LightRequest{
				UserID: userID,
			}
			err := tx.Create(&lightRequest).Error
			if err != nil {
				return err
			}
		}

		lightRequestToLamp := ds.LightRequestToLamp{
			RequestID: lightRequest.ID,
			LampID:    lampID,
		}
		return tx.Create(&lightRequestToLamp).Error
	})
}
