package ds

import "time"

type CalcRequest struct {
	ID                uint64 `gorm:"primaryKey"`
	Status            uint8  `gorm:"not null"`
	UserId            uint64
	User              User `gorm:"foreignKey:UserId"`
	ModeratorId       uint64
	Moderator         User                `gorm:"foreignKey:ModeratorId"`
	MaxTotalPowerW    float64             `gorm:"type:double precision;not null"`
	TotalPowerW       float64             `gorm:"type:double precision;not null"`
	CalcRequestToLamp []CalcRequestToLamp `gorm:"foreignKey:RequestID"`
	CreatedAt         time.Time           `gorm:"not null"`
	FormedAt          time.Time
	ClosedAt          time.Time
}
