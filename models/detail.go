package models

type Detail struct {
	Id          int64     `gorm:"primary_key"`
	TransaksiId int64     `json:"transaksi_id"`
	Transaksi   Transaksi `gorm:"foreign_key:TransaksiId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	MenuId      int64     `json:"menu_id"`
	Menu        Menu      `gorm:"foreign_key:MenuId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Qty         int64     `json:"qty"`
	HargaBeli   float64   `json:"harga_beli"`
}
