package model

type KabKota struct {
	ID         int `gorm:"primaryKey;autoIncrement:false"`
	IdProvinsi int
	Nama       string   `gorm:"type:varchar(255)"`
	Provinsi   Provinsi `gorm:"foreignKey:IdProvinsi"`
}
