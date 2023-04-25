package model

type Provinsi struct {
	ID   int    `gorm:"primaryKey"`
	Nama string `gorm:"type:varchar(255)"`
}
