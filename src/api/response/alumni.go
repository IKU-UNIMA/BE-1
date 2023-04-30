package response

type Alumni struct {
	ID         int    `json:"id"`
	Prodi      string `json:"prodi"`
	Nim        string `json:"nim"`
	Nama       string `json:"nama"`
	Email      string `json:"email"`
	Hp         string `json:"hp"`
	TahunLulus int    `json:"tahun_lulus"`
	Npwp       string `json:"npwp"`
	Nik        int    `json:"nik"`
}
