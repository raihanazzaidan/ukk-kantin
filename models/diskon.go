package models

import (
	"time"
)

type Diskon struct {
	Id               int64     `gorm:"primary_key"`
	NamaDiskon       string    `gorm:"varchar(100)" json:"nama_diskon"`
	PresentaseDiskon string    `gorm:"varchar(100)" json:"presentase_diskon"`
	TanggalAwal      time.Time `json:"tanggal_awal"`
	TanggalAkhir     time.Time `json:"tanggal_akhir"`
	StanId           int64     `json:"stan_id"`
	Stan             Stan      `gorm:"foreign_key:StanId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
