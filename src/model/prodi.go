package model

type Prodi struct {
	ID   int    `gorm:"primaryKey;autoIncrement:false"`
	Nama string `gorm:"type:varchar(255)"`
}
