package repository

import (
	"dia-backend/internal/app/ds"
	"fmt"
	"mime/multipart"

	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

type LampRepository struct {
	db *gorm.DB
}

func NewLampRepository(db *gorm.DB) *LampRepository {
	return &LampRepository{
		db: db,
	}
}

func (r *LampRepository) GetLampByID(id uint64) (*ds.Lamp, error) {
	var lamp ds.Lamp
	err := r.db.Where("is_deleted = false").First(&lamp, id).Error
	if err != nil {
		return nil, err
	}
	return &lamp, nil
}

func (r *LampRepository) GetLamps(titleFilter string) ([]ds.Lamp, error) {
	var lamps []ds.Lamp
	query := r.db.Where("is_deleted = false")

	if titleFilter != "" {
		query = query.Where("title ILIKE ?", "%"+titleFilter+"%")
	}

	err := query.Find(&lamps).Error
	if err != nil {
		return nil, err
	}
	return lamps, nil
}

func (r *LampRepository) CreateLamp(lamp *ds.Lamp) error {
	lamp.IsDeleted = false
	return r.db.Create(lamp).Error
}

func (r *LampRepository) UpdateLamp(id uint64, lampData *ds.Lamp) error {
	return r.db.Model(&ds.Lamp{}).Where("id = ? AND is_deleted = false", id).Updates(map[string]interface{}{
		"title":                lampData.Title,
		"power_w":              lampData.PowerW,
		"luminous_flux_lm":     lampData.LuminousFluxLm,
		"scattering_angle_deg": lampData.ScatteringAngleDeg,
	}).Error
}

func (r *LampRepository) DeleteLamp(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var lamp ds.Lamp
		if err := tx.First(&lamp, id).Error; err != nil {
			return err
		}

		if lamp.ImageURL != "" {
			if err := r.deleteImageFile(lamp.ImageURL); err != nil {
				return err
			}
		}

		return tx.Model(&ds.Lamp{}).Where("id = ?", id).Update("is_deleted", true).Error
	})
}

func (r *LampRepository) AddLampImage(id uint64, fileHeader *multipart.FileHeader) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var lamp ds.Lamp
		if err := tx.Where("is_deleted = false").First(&lamp, id).Error; err != nil {
			return err
		}

		if lamp.ImageURL != "" {
			if err := r.deleteImageFile(lamp.ImageURL); err != nil {
				return err
			}
		}

		fileExt := filepath.Ext(fileHeader.Filename)
		newFileName := fmt.Sprintf("lamp_%d_%d%s", id, time.Now().Unix(), fileExt)
		newFileName = strings.ToLower(newFileName)

		imageURL, err := r.saveImageToMinIO(newFileName)
		if err != nil {
			return err
		}

		return tx.Model(&lamp).Update("image_url", imageURL).Error
	})
}

func (r *LampRepository) AddLampToDraftRequest(lampID uint64, userID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var lamp ds.Lamp
		if err := tx.Where("is_deleted = false").First(&lamp, lampID).Error; err != nil {
			return err
		}

		var lightRequest ds.LightRequest
		err := tx.Where("status = 1 AND user_id = ?", userID).First(&lightRequest).Error
		if err == gorm.ErrRecordNotFound {

			lightRequest = ds.LightRequest{
				Status:    1,
				UserID:    userID,
				CreatedAt: time.Now(),
			}
			if err := tx.Create(&lightRequest).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		lightRequestToLamp := ds.LightRequestToLamp{
			RequestID: lightRequest.ID,
			LampID:    lampID,
			AreaM2:    10.0,
			Number:    1,
		}

		return tx.Create(&lightRequestToLamp).Error
	})
}

func (r *LampRepository) saveImageToMinIO(fileName string) (string, error) {
	return fmt.Sprintf("http://localhost:9000/lamp-images/%s", fileName), nil
}

func (r *LampRepository) deleteImageFile(imageURL string) error {
	if strings.Contains(imageURL, "localhost:9000") {

		fmt.Printf("Image deleted from MinIO: %s\n", imageURL)
	}
	return nil
}
