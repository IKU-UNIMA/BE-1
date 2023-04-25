package model

type Provinsi struct {
	ID   int    `gorm:"primaryKey;autoIncrement:false"`
	Nama string `gorm:"type:varchar(255)"`
}
