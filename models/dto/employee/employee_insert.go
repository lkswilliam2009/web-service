package dto

type EmployeeInsertRequest struct {
	// ===== PEGAWAI =====
	UserID      string     `json:"user_id"`
	NIP         string     `json:"nip"`
	NamaLengkap string     `json:"nama_lengkap"`

	TempatLahir *string    `json:"tempat_lahir"`
	TglLahir    *string  `json:"tgl_lahir"`
	JnsKelamin  *string    `json:"jns_kelamin"`
	Suku        *string    `json:"suku"`
	Agama       *string    `json:"agama"`
	WargaNegara *string    `json:"warga_negara"`
	StatusKawin *string    `json:"status_kawin"`
	GolDarah    *string    `json:"gol_darah"`
	Tinggi      *string    `json:"tinggi"`
	Berat       *string    `json:"berat"`
	Rambut      *string    `json:"rambut"`
	BentukMuka  *string    `json:"bentuk_muka"`
	WarnaKulit  *string    `json:"warna_kulit"`
	CacatTubuh  *string    `json:"cacat_tubuh"`
	Hobby       *string    `json:"hobby"`

	NoKTP      string  `json:"no_ktp"`
	NoNPWP     *string `json:"no_npwp"`
	NoBPJS     *string `json:"no_bpjs"`
	JnsPegawai *string `json:"jns_pegawai"`

	// ===== ALAMAT =====
	AlamatKTP      *string `json:"alamat_ktp"`
	KabKotaID      *string `json:"kabkota_id"`
	KodePos        *string `json:"kodepos"`
	AlamatDomisili *string `json:"alamat_domisili"`
	KabKotaDomID   *string `json:"kabkota_dom_id"`
	KodePosDom     *string `json:"kodepos_dom"`
	NoHP1          *string `json:"no_hp_1"`
	NoHP2          *string `json:"no_hp_2"`
	Email          *string `json:"email"`
}