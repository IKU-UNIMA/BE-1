package request

import (
	"BE-1/src/model"
	"BE-1/src/util"
)

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterAlumni struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required"`
	Nim      string  `json:"nim" validate:"required"`
	Nik      *int    `json:"nik"`
	Npwp     *string `json:"npwp"`
}

type ChangePassword struct {
	PasswordBaru string `json:"password_baru" validate:"required"`
}

func (r *RegisterAlumni) MapRequestToAkun() *model.Akun {
	return &model.Akun{
		Email:    r.Email,
		Password: util.HashPassword(r.Password),
	}
}

func (r *RegisterAlumni) MapRequestToAlumni() *model.Alumni {
	return &model.Alumni{
		Nik:  r.Nik,
		Npwp: r.Npwp,
	}
}
