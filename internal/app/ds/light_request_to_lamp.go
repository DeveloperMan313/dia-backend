package ds

type LightRequestToLamp struct {
	RequestID    uint64       `gorm:"primaryKey" json:"request_id"`
	LampID       uint64       `gorm:"primaryKey" json:"lamp_id"`
	LightRequest LightRequest `gorm:"foreignKey:RequestID;references:ID" json:"-"`
	Lamp         Lamp         `gorm:"foreignKey:LampID;references:ID" json:"lamp"`
	AreaM2       float64      `gorm:"type:double precision;not null;default:10" json:"area_m2"`
	Number       uint64       `gorm:"type:uint" json:"number"`
}
