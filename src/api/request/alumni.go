package request

import (
	"BE-1/src/model"
	"BE-1/src/util"
)

type InsertAlumni struct {
	IdProdi    int    `json:"id_prodi" validate:"required"`
	Nim        string `json:"nim" validate:"required"`
	Nama       string `json:"nama" validate:"required"`
	Hp         string `json:"hp" validate:"required"`
	TahunLulus uint   `json:"tahun_lulus" validate:"required"`
}

type EditAlumni struct {
	IdProdi    int     `json:"id_prodi" validate:"required"`
	Nim        string  `json:"nim" validate:"required"`
	Nama       string  `json:"nama" validate:"required"`
	Hp         string  `json:"hp" validate:"required"`
	Email      *string `json:"email"`
	TahunLulus uint    `json:"tahun_lulus" validate:"required"`
	Npwp       *string `json:"npwp"`
	Nik        *int    `json:"nik"`
}

func (r *InsertAlumni) MapRequest() *model.Alumni {
	return &model.Alumni{
		IdProdi:    r.IdProdi,
		KodePt:     util.KODE_PT,
		Nim:        r.Nim,
		Nama:       r.Nama,
		Hp:         r.Hp,
		TahunLulus: r.TahunLulus,
	}
}

func (r *EditAlumni) MapRequest() *model.Alumni {
	return &model.Alumni{
		IdProdi:    r.IdProdi,
		Nim:        r.Nim,
		Nama:       r.Nama,
		Hp:         r.Hp,
		Email:      r.Email,
		TahunLulus: r.TahunLulus,
		Npwp:       r.Npwp,
		Nik:        r.Nik,
	}
}
