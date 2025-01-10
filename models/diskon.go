package models

import (
	"time"
)

type Diskon struct {
	Id               int64  `gorm:"primary_key"`
	NamaDiskon       string `gorm:"varchar(100)" json:"nama_diskon"`
	PresentaseDiskon string `gorm:"varchar(100)" json:"presentase_diskon"`
	TanggalAwal 	 time.Time `json:"tanggal_awal"`
	TanggalAkhir 	 time.Time `json:"tanggal_akhir"`
}
