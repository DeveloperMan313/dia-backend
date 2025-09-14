package repository

import (
	"dia-backend/internal/app/ds"

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
	err := r.db.Where("id = ?", id).Take(&lamp).Error
	if err != nil {
		return nil, err
	}

	return &lamp, nil
}

func (r *LampRepository) GetLamps() ([]ds.Lamp, error) {
	var lamps []ds.Lamp
	err := r.db.Find(&lamps).Error
	if err != nil {
		return nil, err
	}

	return lamps, nil
}

func (r *LampRepository) GetLampsByTitle(title string) ([]ds.Lamp, error) {
	var lamps []ds.Lamp
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&lamps).Error
	if err != nil {
		return nil, err
	}

	return lamps, nil
}
