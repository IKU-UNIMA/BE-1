package response

type Alumni struct {
	ID         int    `json:"id"`
	IdProdi    int    `json:"-"`
	KodePt     string `json:"kode_pt"`
	Prodi      Prodi  `gorm:"foreignKey:IdProdi" json:"prodi"`
	Nim        string `json:"nim"`
	Nama       string `json:"nama"`
	Email      string `json:"email"`
	Hp         string `json:"hp"`
	TahunLulus int    `json:"tahun_lulus"`
	Npwp       string `json:"npwp"`
	Nik        int    `json:"nik"`
}
