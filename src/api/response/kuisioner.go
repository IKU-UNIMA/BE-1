package response

type (
	Kuisioner struct {
		ID       int    `json:"id"`
		IdAlumni string `json:"-"`
		Alumni   Alumni `gorm:"foreignKey:IdAlumni" json:"alumni"`
		Status   bool   `json:"status"`
	}

	DetailKuisioner struct {
		ID       int    `json:"id"`
		IdAlumni string `json:"-"`
		Alumni   Alumni `gorm:"foreignKey:IdAlumni" json:"alumni"`
		F8       int8   `json:"f8"`
		F504     int8   `json:"f504"`
		F502     int8   `json:"f502"`
		F505     int    `json:"f505"`
		F5a1     int    `json:"f5a1"`
		F5a2     int    `json:"f5a2"`
		F1101    int8   `json:"f1101"`
		F1102    string `json:"f1102"`
		F5b      string `json:"f5b"`
		F5c      int8   `json:"f5c"`
		F5d      int8   `json:"f5d"`
		F18a     int8   `json:"f18a"`
		F18b     string `json:"f18b"`
		F18c     string `json:"f18c"`
		F18d     string `json:"f18d"`
		F1201    int8   `json:"f1201"`
		F1202    string `json:"f1202"`
		F14      int8   `json:"f14"`
		F15      int8   `json:"f15"`
		F1761    int8   `json:"f1761"`
		F1762    int8   `json:"f1762"`
		F1763    int8   `json:"f1763"`
		F1764    int8   `json:"f1764"`
		F1765    int8   `json:"f1765"`
		F1766    int8   `json:"f1766"`
		F1767    int8   `json:"f1767"`
		F1768    int8   `json:"f1768"`
		F1769    int8   `json:"f1769"`
		F1770    int8   `json:"f1770"`
		F1771    int8   `json:"f1771"`
		F1772    int8   `json:"f1772"`
		F1773    int8   `json:"f1773"`
		F1774    int8   `json:"f1774"`
		F21      int8   `json:"f21"`
		F22      int8   `json:"f22"`
		F23      int8   `json:"f23"`
		F24      int8   `json:"f24"`
		F25      int8   `json:"f25"`
		F26      int8   `json:"f26"`
		F27      int8   `json:"f27"`
		F301     int8   `json:"f28"`
		F302     int    `json:"f302"`
		F303     int    `json:"f303"`
		F401     int8   `json:"f401"`
		F402     int8   `json:"f402"`
		F403     int8   `json:"f403"`
		F404     int8   `json:"f404"`
		F405     int8   `json:"f405"`
		F406     int8   `json:"f406"`
		F407     int8   `json:"f407"`
		F408     int8   `json:"f408"`
		F409     int8   `json:"f409"`
		F410     int8   `json:"f410"`
		F411     int8   `json:"f411"`
		F412     int8   `json:"f412"`
		F413     int8   `json:"f413"`
		F414     int8   `json:"f414"`
		F415     int8   `json:"f415"`
		F416     string `json:"f4016"`
		F6       int    `json:"f6"`
		F7       int    `json:"f7"`
		F7a      int    `json:"f7a"`
		F1001    int8   `json:"f1001"`
		F1002    string `json:"f1002"`
		F1601    int8   `json:"f1601"`
		F1602    int8   `json:"f1602"`
		F1603    int8   `json:"f1603"`
		F1604    int8   `json:"f1604"`
		F1605    int8   `json:"f1605"`
		F1606    int8   `json:"f1606"`
		F1607    int8   `json:"f1607"`
		F1608    int8   `json:"f1608"`
		F1609    int8   `json:"f1609"`
		F1610    int8   `json:"f1610"`
		F1611    int8   `json:"f1611"`
		F1612    int8   `json:"f1612"`
		F1613    int8   `json:"f1613"`
		F1614    string `json:"f1614"`
	}
)
