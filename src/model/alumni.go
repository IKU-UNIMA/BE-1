package model

type Alumni struct {
	ID         int `gorm:"primaryKey"`
	IdProdi    uint
	KodePt     string `gorm:"type:varchar(15)"`
	Nim        string `gorm:"type:varchar(10);unique"`
	Nama       string `gorm:"type:varchar(255)"`
	Hp         string `gorm:"type:varchar(20);unique"`
	TahunLulus uint   `gorm:"type:smallint"`
	Npwp       string `gorm:"type:varchar(255);unique"`
	Prodi      Prodi  `gorm:"foreignKey:IdProdi"`
}
