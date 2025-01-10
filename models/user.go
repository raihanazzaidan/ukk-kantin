package models

import (
	"time"
)

type User struct {
	Id        int64     `gorm:"primary_key"`
	Username  string    `gorm:"type:varchar(100)" json:"username"`
	Password  string    `gorm:"type:varchar(100)" json:"password"`
	RoleId    int64     `json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
}

type Role struct {
	Id   int64  `gorm:"primary_key"`
	Role string `gorm:"type:varchar(20)" json:"role"`
}
