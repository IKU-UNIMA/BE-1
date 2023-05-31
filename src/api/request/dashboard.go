package request

import (
	"BE-1/src/model"
	"BE-1/src/util"
)

type Target struct {
	Target float32 `json:"target" validate:"required"`
	Tahun  int     `json:"tahun" validate:"required"`
}

func (r *Target) MapRequest() *model.Target {
	return &model.Target{
		Bagian: util.IKU1,
		Target: r.Target,
		Tahun:  r.Tahun,
	}
}
