package ds

type LightRequestToLamp struct {
	RequestID    uint64       `gorm:"primaryKey" json:"-"`
	LampID       uint64       `gorm:"primaryKey" json:"-"`
	LightRequest LightRequest `gorm:"foreignKey:RequestID;references:ID" json:"-"`
	Lamp         Lamp         `gorm:"foreignKey:LampID;references:ID" json:"lamp"`
	AreaM2       float64      `gorm:"type:double precision;default:null" json:"area_m2"`
	Number       uint64       `gorm:"type:uint;default:null" json:"number"`
}
