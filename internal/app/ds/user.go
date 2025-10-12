package ds

import "dia-backend/internal/app/role"

type User struct {
	ID       uint64    `gorm:"primaryKey" json:"-"`
	Username string    `gorm:"type:varchar(50);not null;unique" json:"username"`
	Password string    `gorm:"column:passwrd;type:varchar(255);not null" json:"-"`
	Role     role.Role `gorm:"type:int;not null;default:0" json:"role"`
}
