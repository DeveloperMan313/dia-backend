package repository

import (
	"context"
	"dia-backend/internal/app/ds"
	"errors"
	"fmt"
	"mime/multipart"
	"os"

	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LampRepository struct {
	db          *gorm.DB
	minioClient *minio.Client
}

func NewLampRepository(db *gorm.DB, minioClient *minio.Client) *LampRepository {
	return &LampRepository{
		db:          db,
		minioClient: minioClient,
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

		imageURL, err := r.saveLampImageToMinIO(newFileName, fileHeader)
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

const lampImagesBucket = "lamp-images"

func (r *LampRepository) saveLampImageToMinIO(fileName string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileSize := fileHeader.Size

	contentType := "application/octet-stream"
	if strings.HasSuffix(strings.ToLower(fileName), ".jpg") || strings.HasSuffix(strings.ToLower(fileName), ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(strings.ToLower(fileName), ".gif") {
		contentType = "image/gif"
	}

	_, err = r.minioClient.PutObject(context.Background(), lampImagesBucket, fileName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s/%s/%s", os.Getenv("MINIO_HOST"), os.Getenv("MINIO_SERVER_PORT"), lampImagesBucket, fileName), nil
}

func (r *LampRepository) deleteImageFile(imageURL string) error {
	minioOrigin := os.Getenv("MINIO_HOST") + ":" + os.Getenv("MINIO_SERVER_PORT")
	if strings.Contains(imageURL, minioOrigin) {
		parts := strings.Split(imageURL, "/")
		if len(parts) > 0 {
			fileName := parts[len(parts)-1]
			err := r.minioClient.RemoveObject(context.Background(), lampImagesBucket, fileName, minio.RemoveObjectOptions{})
			if err != nil {
				return err
			}
			logrus.Printf("Image deleted from MinIO: %s\n", imageURL)
			return nil
		}
	}
	return errors.New("could not delete image file")
}
