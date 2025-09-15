package ds

type CalcRequestToLamp struct {
	RequestID   uint64      `gorm:"primaryKey"`
	LampID      uint64      `gorm:"primaryKey"`
	CalcRequest CalcRequest `gorm:"foreignKey:RequestID;references:ID"`
	Lamp        Lamp        `gorm:"foreignKey:LampID;references:ID"`
	AreaM2      float64     `gorm:"type:double precision;not null;default:10"`
	Number      uint64      `gorm:"type:uint"`
}
