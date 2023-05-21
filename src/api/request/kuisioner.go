package request

import "BE-1/src/model"

type Kuisioner struct {
	Nim   string  `json:"nim" validate:"required"`
	Nik   *int    `json:"nik" validate:"required"`
	Npwp  *string `json:"npwp"`
	Email *string `json:"email" validate:"required,email"`
	Hp    string  `json:"hp" validate:"required"`
	F8    int8    `json:"f8" validate:"required,min=1,max=5"`
	F504  int8    `json:"f504" validate:"min=1,max=2"`
	F502  int8    `json:"f502" validate:"min=0"`
	F505  int32   `json:"505" validate:"min=0"`
	F5a1  int32   `json:"5a1"`
	F5a2  int32   `json:"f5a2"`
	F1101 int8    `json:"f1101" validate:"min=1,max=7"`
	F1102 string  `json:"f1102"`
	F5b   string  `json:"f5b"`
	F5c   int8    `json:"f5c" validate:"min=1,max=4"`
	F5d   int8    `json:"f5d" validate:"min=1,max=3"`
	F18a  int8    `json:"f18a" validate:"min=1,max=2"`
	F18b  string  `json:"f18b"`
	F18c  string  `json:"f18c"`
	F18d  string  `json:"f18d"`
	F1201 int8    `json:"f1201" validate:"required,min=1,max=7"`
	F1202 string  `json:"f1202"`
	F14   int8    `json:"f14" validate:"required,min=1,max=5"`
	F15   int8    `json:"f15" validate:"required,min=1,max=4"`
	F1761 int8    `json:"f1761" validate:"required,min=1,max=5"`
	F1762 int8    `json:"f1762" validate:"required,min=1,max=5"`
	F1763 int8    `json:"f1763" validate:"required,min=1,max=5"`
	F1764 int8    `json:"f1764" validate:"required,min=1,max=5"`
	F1765 int8    `json:"f1765" validate:"required,min=1,max=5"`
	F1766 int8    `json:"f1766" validate:"required,min=1,max=5"`
	F1767 int8    `json:"f1767" validate:"required,min=1,max=5"`
	F1768 int8    `json:"f1768" validate:"required,min=1,max=5"`
	F1769 int8    `json:"f1769" validate:"required,min=1,max=5"`
	F1770 int8    `json:"f1770" validate:"required,min=1,max=5"`
	F1771 int8    `json:"f1771" validate:"required,min=1,max=5"`
	F1772 int8    `json:"f1772" validate:"required,min=1,max=5"`
	F1773 int8    `json:"f1773" validate:"required,min=1,max=5"`
	F1774 int8    `json:"f1774" validate:"required,min=1,max=5"`
	F21   int8    `json:"f21" validate:"min=1,max=5"`
	F22   int8    `json:"f22" validate:"min=1,max=5"`
	F23   int8    `json:"f23" validate:"min=1,max=5"`
	F24   int8    `json:"f24" validate:"min=1,max=5"`
	F25   int8    `json:"f25" validate:"min=1,max=5"`
	F26   int8    `json:"f26" validate:"min=1,max=5"`
	F27   int8    `json:"f27" validate:"min=1,max=5"`
	F301  int8    `json:"f301" validate:"min=1,max=3"`
	F302  int8    `json:"f302" validate:"min=0"`
	F303  int8    `json:"f303" validate:"min=0"`
	F401  int8    `json:"f401" validate:"min=0,max=1"`
	F402  int8    `json:"f402" validate:"min=0,max=1"`
	F403  int8    `json:"f403" validate:"min=0,max=1"`
	F404  int8    `json:"f404" validate:"min=0,max=1"`
	F405  int8    `json:"f405" validate:"min=0,max=1"`
	F406  int8    `json:"f406" validate:"min=0,max=1"`
	F407  int8    `json:"f407" validate:"min=0,max=1"`
	F408  int8    `json:"f408" validate:"min=0,max=1"`
	F409  int8    `json:"f409" validate:"min=0,max=1"`
	F410  int8    `json:"f410" validate:"min=0,max=1"`
	F411  int8    `json:"f411" validate:"min=0,max=1"`
	F412  int8    `json:"f412" validate:"min=0,max=1"`
	F413  int8    `json:"f413" validate:"min=0,max=1"`
	F414  int8    `json:"f414" validate:"min=0,max=1"`
	F415  int8    `json:"f415" validate:"min=0,max=1"`
	F416  string  `json:"f4016"`
	F6    int16   `json:"f6" validate:"min=0"`
	F7    int16   `json:"f7" validate:"min=0"`
	F7a   int16   `json:"f7a" validate:"min=0"`
	F1001 int8    `json:"f1001" validate:"min=1,max=5"`
	F1002 string  `json:"f1002"`
	F1601 int8    `json:"f1601" validate:"min=0,max=1"`
	F1602 int8    `json:"f1602" validate:"min=0,max=1"`
	F1603 int8    `json:"f1603" validate:"min=0,max=1"`
	F1604 int8    `json:"f1604" validate:"min=0,max=1"`
	F1605 int8    `json:"f1605" validate:"min=0,max=1"`
	F1606 int8    `json:"f1606" validate:"min=0,max=1"`
	F1607 int8    `json:"f1607" validate:"min=0,max=1"`
	F1608 int8    `json:"f1608" validate:"min=0,max=1"`
	F1609 int8    `json:"f1609" validate:"min=0,max=1"`
	F1610 int8    `json:"f1610" validate:"min=0,max=1"`
	F1611 int8    `json:"f1611" validate:"min=0,max=1"`
	F1612 int8    `json:"f1612" validate:"min=0,max=1"`
	F1613 int8    `json:"f1613" validate:"min=0,max=1"`
	F1614 string  `json:"f1614"`
}

func (r *Kuisioner) MapAlumniData() *model.Alumni {
	return &model.Alumni{
		Nik:   r.Nik,
		Npwp:  r.Npwp,
		Email: r.Email,
		Hp:    r.Hp,
	}
}

func (r *Kuisioner) MapRequest(idAlumni int) *model.Kuisioner {
	parseBool := func(v int8) bool {
		return v == 1
	}

	return &model.Kuisioner{
		IdAlumni: idAlumni,
		F8:       r.F8,
		F504:     r.F504,
		F502:     r.F502,
		F505:     r.F505,
		F5a1:     r.F5a1,
		F5a2:     r.F5a2,
		F1101:    r.F1101,
		F1102:    r.F1102,
		F5b:      r.F5b,
		F5c:      r.F5c,
		F5d:      r.F5d,
		F18a:     r.F18a,
		F18b:     r.F18b,
		F18c:     r.F18c,
		F18d:     r.F18d,
		F1201:    r.F1201,
		F1202:    r.F1202,
		F14:      r.F14,
		F15:      r.F15,
		F1761:    r.F1761,
		F1762:    r.F1762,
		F1763:    r.F1763,
		F1764:    r.F1764,
		F1765:    r.F1765,
		F1766:    r.F1766,
		F1767:    r.F1767,
		F1768:    r.F1768,
		F1769:    r.F1769,
		F1770:    r.F1770,
		F1771:    r.F1771,
		F1772:    r.F1772,
		F1773:    r.F1773,
		F1774:    r.F1774,
		F21:      r.F21,
		F22:      r.F22,
		F23:      r.F23,
		F24:      r.F24,
		F25:      r.F25,
		F26:      r.F26,
		F27:      r.F27,
		F301:     r.F301,
		F302:     r.F302,
		F303:     r.F303,
		F401:     parseBool(r.F401),
		F402:     parseBool(r.F402),
		F403:     parseBool(r.F403),
		F404:     parseBool(r.F404),
		F405:     parseBool(r.F405),
		F406:     parseBool(r.F406),
		F407:     parseBool(r.F407),
		F408:     parseBool(r.F408),
		F409:     parseBool(r.F409),
		F410:     parseBool(r.F410),
		F411:     parseBool(r.F411),
		F412:     parseBool(r.F412),
		F413:     parseBool(r.F413),
		F414:     parseBool(r.F414),
		F415:     parseBool(r.F415),
		F416:     r.F416,
		F6:       r.F6,
		F7:       r.F7,
		F7a:      r.F7a,
		F1001:    r.F1001,
		F1002:    r.F1102,
		F1601:    parseBool(r.F1601),
		F1602:    parseBool(r.F1602),
		F1603:    parseBool(r.F1603),
		F1604:    parseBool(r.F1604),
		F1605:    parseBool(r.F1605),
		F1606:    parseBool(r.F1606),
		F1607:    parseBool(r.F1607),
		F1608:    parseBool(r.F1608),
		F1609:    parseBool(r.F1609),
		F1610:    parseBool(r.F1610),
		F1611:    parseBool(r.F1611),
		F1612:    parseBool(r.F1612),
		F1613:    parseBool(r.F1613),
		F1614:    r.F1614,
	}
}
