package repository

import (
	"context"
	"dia-backend/internal/app/config"
	"dia-backend/internal/app/dsn"
	redisClient "dia-backend/internal/app/redis"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db           *gorm.DB
	Lamp         *LampRepository
	LightRequest *LightRequestRepository
	User         *UserRepository
	Redis        *redisClient.Client
	Config       *config.Config
}

func NewRepository() (*Repository, error) {
	cfg := config.LoadConfig()

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	minioClient, err := InitMinIOClient()
	if err != nil {
		return nil, err
	}

	redis, err := redisClient.New(cfg.Redis)
	if err != nil {
		logrus.Warnf("Failed to connect to Redis: %v. JWT blacklist will not work.", err)
	}

	return &Repository{
		db:           db,
		Lamp:         NewLampRepository(db, minioClient),
		LightRequest: NewLightRequestRepository(db),
		User:         NewUserRepository(db),
		Redis:        redis,
		Config:       cfg,
	}, nil
}

func CloseDBConn(r *Repository) {
	dbInstance, _ := r.db.DB()
	_ = dbInstance.Close()
	if r.Redis != nil {
		_ = r.Redis.Close()
	}
}

func InitMinIOClient() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_HOST") + ":" + os.Getenv("MINIO_SERVER_PORT")
	accessKeyID := os.Getenv("MINIO_ROOT_USER")
	secretAccessKey := os.Getenv("MINIO_ROOT_PASSWORD")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %v", err)
	}

	ctx := context.Background()

	exists, err := minioClient.BucketExists(ctx, lampImagesBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %v", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, lampImagesBucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %v", err)
		}
		logrus.Printf("Bucket '%s' created successfully\n", lampImagesBucket)
	}

	return minioClient, nil
}
