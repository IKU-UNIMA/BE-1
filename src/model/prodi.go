package model

type Prodi struct {
	ID         int `gorm:"primaryKey"`
	IdFakultas int
	KodeProdi  int
	Nama       string `gorm:"type:varchar(255)"`
	Jenjang    string `gorm:"type:varchar(60)"`
}
