package user

import "time"

type Entity struct {
	UserID    uint64    `gorm:"column:user_id;primaryKey;autoIncrement:true"`
	Email     string    `gorm:"column:email;size:256;unique"`
	Password  string    `gorm:"column:password;size:256"`
	Username  string    `gorm:"column:username;size:256;unique"`
	Fullname  string    `gorm:"column:fullname;size:256"`
	Phone     string    `gorm:"column:phone;size:256"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:current_timestamp;autoUpdateTime"`
}

func (a *Entity) TableName() string {
	return "user"
}
