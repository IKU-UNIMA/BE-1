package response

type (
	Dashboard struct {
		Target      string                       `json:"target"`
		Total       int                          `json:"total"`
		TotalAlumni int                          `json:"total_alumni"`
		Pencapaian  string                       `json:"pencapaian"`
		Detail      []DashboardDetailPerFakultas `json:"detail"`
	}

	DashboardDetailPerFakultas struct {
		ID                    int    `json:"id"`
		Fakultas              string `json:"fakultas"`
		JumlahAlumni          int    `json:"jumlah_alumni"`
		Bekerja               int    `json:"bekerja"`
		Wiraswasta            int    `json:"wiraswasta"`
		MelanjutkanPendidikan int    `json:"melanjutkan_pendidikan"`
		JumlahResponden       int    `json:"jumlah_responden"`
		Persentase            string `json:"persentase"`
	}

	DashboardPerProdi struct {
		Fakultas    string                    `json:"fakultas"`
		Total       int                       `json:"total"`
		TotalAlumni int                       `json:"total_alumni"`
		Pencapaian  string                    `json:"pencapaian"`
		Detail      []DashboardDetailPerProdi `json:"detail"`
	}

	DashboardDetailPerProdi struct {
		Prodi                 string `json:"prodi"`
		JumlahAlumni          int    `json:"jumlah_alumni"`
		Bekerja               int    `json:"bekerja"`
		Wiraswasta            int    `json:"wiraswasta"`
		MelanjutkanPendidikan int    `json:"melanjutkan_pendidikan"`
		JumlahResponden       int    `json:"jumlah_responden"`
		Persentase            string `json:"persentase"`
	}
)
