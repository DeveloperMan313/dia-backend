package ds

type Lamp struct {
	ID                 uint64               `gorm:"primaryKey"`
	Title              string               `gorm:"type:varchar(100);not null"`
	PowerW             float64              `gorm:"type:double precision;not null"`
	LuminousFluxLm     float64              `gorm:"type:double precision;not null"`
	ScatteringAngleDeg float64              `gorm:"type:double precision;not null"`
	IsDeleted          bool                 `gorm:"boolean;not null"`
	ImageURL           string               `gorm:"type:varchar(100)"`
	LightRequestToLamp []LightRequestToLamp `gorm:"foreignKey:LampID"`
}
