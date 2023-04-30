package model

type Alumni struct {
	ID         int `gorm:"primaryKey"`
	IdProdi    int
	KodePt     string  `gorm:"type:varchar(15)"`
	Nim        string  `gorm:"type:varchar(10);unique"`
	Nama       string  `gorm:"type:varchar(255)"`
	Hp         string  `gorm:"type:varchar(20)"`
	TahunLulus uint    `gorm:"type:smallint"`
	Npwp       *string `gorm:"type:varchar(255);unique"`
	Nik        *int    `gorm:"unique"`
	Prodi      Prodi   `gorm:"foreignKey:IdProdi"`
	Akun       Akun    `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
}
