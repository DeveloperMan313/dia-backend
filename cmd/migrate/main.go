package main

import (
	"dia-backend/internal/app/ds"
	"dia-backend/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load("deploy/.env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&ds.Lamp{},
		&ds.CalcRequest{},
		&ds.CalcRequestToLamp{},
		&ds.User{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
