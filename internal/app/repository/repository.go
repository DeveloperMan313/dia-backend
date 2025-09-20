package repository

import (
	"dia-backend/internal/app/dsn"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	Lamp         *LampRepository
	LightRequest *LightRequestRepository
}

func NewRepository() (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		Lamp:         NewLampRepository(db),
		LightRequest: NewLightRequestRepository(db),
	}, nil
}
