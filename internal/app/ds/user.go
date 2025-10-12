package ds

type User struct {
	ID       uint64 `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type:varchar(50);not null;unique" json:"username"`
	Passwrd  string `gorm:"type:varchar(50);not null" json:"-"`
	IsMod    bool   `gorm:"type:boolean;not null" json:"is_mod"`
}
