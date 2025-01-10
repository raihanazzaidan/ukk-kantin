package models

import (
	"time"
)

type Transaksi struct {
	Id       int64     `gorm:"primary_key"`
	Tanggal  time.Time `gorm:"type:timestamp;default:current_timestamp" json:"tanggal"`
	StanId   int64     `json:"stan_id"`
	Stan     Stan      `gorm:"foreign_key:StanId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	SiswaId  int64     `json:"siswa_id"`
	Siswa    Siswa     `gorm:"foreign_key:SiswaId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	StatusId int64     `json:"status"`
	Status   Status    `gorm:"foreign_key:StatusId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Status struct {
	Id     int64  `gorm:"primary_key"`
	Status string `gorm:"type:varchar(50)" json:"status"`
}
