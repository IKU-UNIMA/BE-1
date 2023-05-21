package model

type Prodi struct {
	ID        int `gorm:"primaryKey"`
	KodeProdi int
	Nama      string `gorm:"type:varchar(255)"`
}
