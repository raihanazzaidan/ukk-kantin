package models

import (
	"time"
)

type Siswa struct {
	Id        int64     `gorm:"primary_key"`
	NamaSiswa string    `gorm:"type:varchar(100)" json:"nama_siswa"`
	Alamat    string    `json:"alamat"`
	Telp      string    `gorm:"type:varchar(20)" json:"telp"`
	UserId    int64     `json:"user_id"`
	User      User      `gorm:"foreignKey:UserId;references:Id;	constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Foto      string    `json:"foto"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:current_timestamp on update current_timestamp"`
}
