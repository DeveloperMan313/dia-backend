package ds

type Lamp struct {
	ID                 uint64               `gorm:"primaryKey" json:"id"`
	Title              string               `gorm:"type:varchar(100);not null" json:"title"`
	PowerW             float64              `gorm:"type:double precision;not null" json:"power_w"`
	LuminousFluxLm     float64              `gorm:"type:double precision;not null" json:"luminous_flux_lm"`
	ScatteringAngleDeg float64              `gorm:"type:double precision;not null" json:"scattering_angle_deg"`
	IsDeleted          bool                 `gorm:"boolean;not null" json:"-"`
	ImageURL           string               `gorm:"type:varchar(100)" json:"image_url"`
	LightRequestToLamp []LightRequestToLamp `gorm:"foreignKey:LampID" json:"-"`
}
