package models

type MenuDiskon struct {
	Id       int64  `gorm:"primary_key"`
	MenuId   int64  `json:"menu_id"`
	Menu     Menu   `gorm:"foreignKey:MenuId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DiskonId int64  `json:"diskon_id"`
	Diskon   Diskon `gorm:"foreignKey:DiskonId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StanId   int64  `json:"stan_id"`
	Stan     Stan   `gorm:"foreignKey:StanId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
