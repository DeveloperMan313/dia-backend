package ds

import "time"

type LightRequest struct {
	ID                 uint64               `gorm:"primaryKey" json:"id"`
	Status             uint8                `gorm:"not null;default:1" json:"status"` // 1 - draft, 2 - deleted, 3 - pending, 4 - resolved, 5 - rejected
	UserID             uint64               `gorm:"not null" json:"user_id"`
	User               User                 `gorm:"foreignKey:UserID;references:ID" json:"-"`
	ModeratorId        uint64               `gorm:"default:null" json:"-"`
	Moderator          User                 `gorm:"foreignKey:ModeratorId;references:ID" json:"-"`
	MaxTotalPowerW     float64              `gorm:"type:double precision;default:null" json:"max_total_power_w"`
	LightRequestToLamp []LightRequestToLamp `gorm:"foreignKey:RequestID" json:"light_request_to_lamp"`
	CreatedAt          time.Time            `gorm:"not null;default:now()" json:"created_at"`
	FormedAt           time.Time            `gorm:"default:null" json:"formed_at"`
	ClosedAt           time.Time            `gorm:"default:null" json:"closed_at"`
	CalculatedCnt      uint64               `json:"calculated_cnt"`
}
