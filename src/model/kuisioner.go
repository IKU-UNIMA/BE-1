package model

import "time"

type Kuisioner struct {
	ID        int    `gorm:"primaryKey"`
	IdAlumni  int    `gorm:"unique"`
	Alumni    Alumni `gorm:"foreignKey:IdAlumni;constraint:OnDelete:CASCADE"`
	F8        int8
	F504      int8
	F502      int8
	F505      int32
	F5a1      int32
	F5a2      int32
	F1101     int8
	F1102     string `gorm:"type:varchar(255)"`
	F5b       string `gorm:"type:varchar(255)"`
	F5c       int8
	F5d       int8
	F18a      int8
	F18b      string `gorm:"type:varchar(255)"`
	F18c      string `gorm:"type:varchar(255)"`
	F18d      string `gorm:"type:varchar(255)"`
	F1201     int8
	F1202     string `gorm:"type:varchar(255)"`
	F14       int8
	F15       int8
	F1761     int8
	F1762     int8
	F1763     int8
	F1764     int8
	F1765     int8
	F1766     int8
	F1767     int8
	F1768     int8
	F1769     int8
	F1770     int8
	F1771     int8
	F1772     int8
	F1773     int8
	F1774     int8
	F21       int8
	F22       int8
	F23       int8
	F24       int8
	F25       int8
	F26       int8
	F27       int8
	F301      int8
	F302      int8
	F303      int8
	F401      bool
	F402      bool
	F403      bool
	F404      bool
	F405      bool
	F406      bool
	F407      bool
	F408      bool
	F409      bool
	F410      bool
	F411      bool
	F412      bool
	F413      bool
	F414      bool
	F415      bool
	F416      string `gorm:"type:varchar(255)"`
	F6        int16
	F7        int16
	F7a       int16
	F1001     int8
	F1002     string `gorm:"type:varchar(255)"`
	F1601     bool
	F1602     bool
	F1603     bool
	F1604     bool
	F1605     bool
	F1606     bool
	F1607     bool
	F1608     bool
	F1609     bool
	F1610     bool
	F1611     bool
	F1612     bool
	F1613     bool
	F1614     string `gorm:"type:varchar(255)"`
	Status    *bool
	CreatedAt time.Time
}
