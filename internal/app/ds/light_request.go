package ds

import "time"

type LightRequest struct {
	ID                 uint64               `gorm:"primaryKey"`
	Status             uint8                `gorm:"not null;default:1"` // 1 - draft, 2 - deleted, 3 - pending, 4 - resolved, 5 - rejected
	UserID             uint64               `gorm:"not null"`
	User               User                 `gorm:"foreignKey:UserID;references:ID"`
	ModeratorId        uint64               `gorm:"default:null"`
	Moderator          User                 `gorm:"foreignKey:ModeratorId;references:ID"`
	MaxTotalPowerW     float64              `gorm:"type:double precision;not null;default:1000"`
	TotalPowerW        float64              `gorm:"type:double precision"`
	LightRequestToLamp []LightRequestToLamp `gorm:"foreignKey:RequestID"`
	CreatedAt          time.Time            `gorm:"not null;default:now()"`
	FormedAt           time.Time            `gorm:"default:null"`
	ClosedAt           time.Time            `gorm:"default:null"`
}
