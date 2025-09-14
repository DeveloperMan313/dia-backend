package ds

type CalcRequestToLamp struct {
	RequestID   uint64
	CalcRequest CalcRequest `gorm:"foreignKey:RequestID"`
	LampID      uint64
	Lamp        Lamp    `gorm:"foreignKey:LampID"`
	AreaM2      float64 `gorm:"type:double precision;not null"`
	Number      uint64  `gorm:"type:uint"`
}
