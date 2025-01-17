package models

type Menu struct {
	Id          int64     `gorm:"primary_key"`
	NamaMakanan string    `gorm:"type:varchar(100)" json:"nama_makanan"`
	Harga       float64   `json:"harga"`
	JenisId     int64     `json:"jenis_id"`
	Jenis       Jenis     `gorm:"foreign_key:JenisId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Foto        string    `gorm:"type:varchar(255)" json:"foto"`
	Deskripsi   string    `json:"deskripsi"`
	StanId      int64     `json:"stan_id"`
	Stan        Stan      `gorm:"foreign_key:StanId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Jenis struct {
	Id    int64  `gorm:"primary_key"`
	Jenis string `gorm:"type:varchar(20)" json:"jenis"`
}
