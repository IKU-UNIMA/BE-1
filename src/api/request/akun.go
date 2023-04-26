package request

import (
	"BE-1/src/model"
	"BE-1/src/util"
)

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CheckNIM struct {
	Nim string `json:"nim" validate:"required"`
}

type RegisterAlumni struct {
	Nim      string `json:"nim" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ChangePassword struct {
	PasswordBaru string `json:"password_baru" validate:"required"`
}

func (r *RegisterAlumni) MapRequest() *model.Akun {
	return &model.Akun{
		Email:    r.Email,
		Password: util.HashPassword(r.Password),
	}
}
