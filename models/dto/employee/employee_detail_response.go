package dto

import (
	"time"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type EmployeeDetailResponse struct {
	PegawaiID uuid.UUID `db:"pegawaiID" json:"pegawai_id"`
	UserID    *uuid.UUID `db:"userid" json:"user_id"`
	BranchID  *uuid.UUID `db:"branchid" json:"branch_id"`

	NIP          string `db:"NIP" json:"nip"`
	NamaLengkap  string `db:"nama_lengkap" json:"nama_lengkap"`
	TempatLahir  *string `db:"tempat_lahir" json:"tempat_lahir"`
	TglLahir     *time.Time `db:"tgl_lahir" json:"tgl_lahir"`
	JnsKelamin   *string `db:"jns_kelamin" json:"jns_kelamin"`
	Suku         *string `db:"suku" json:"suku"`
	Agama        *string `db:"agama" json:"agama"`
	WargaNegara  *string `db:"warga_negara" json:"warga_negara"`
	StatusKawin  *string `db:"status_kawin" json:"status_kawin"`
	GolDarah     *string `db:"gol_darah" json:"gol_darah"`

	Tinggi *int `db:"tinggi" json:"tinggi"`
	Berat  *int `db:"berat" json:"berat"`

	Rambut      *string `db:"rambut" json:"rambut"`
	BentukMuka  *string `db:"bentuk_muka" json:"bentuk_muka"`
	WarnaKulit  *string `db:"warna_kulit" json:"warna_kulit"`
	CacatTubuh  *string `db:"cacat_tubuh" json:"cacat_tubuh"`
	Hobby       *string `db:"hobby" json:"hobby"`

	NoKTP  *string `db:"no_KTP" json:"no_ktp"`
	NoNPWP *string `db:"no_NPWP" json:"no_npwp"`
	NoBPJS *string `db:"no_BPJS" json:"no_bpjs"`

	AlamatKTP       *string `db:"alamat_ktp" json:"alamat_ktp"`
	KabKotaID       *int `db:"kabkotaid" json:"kabkota_id"`
	KabKota         *string `db:"kab_kota" json:"kab_kota"`
	ProvID          *int `db:"provid" json:"prov_id"`
	Provinsi        *string `db:"provinsi" json:"provinsi"`
	KodePos         *string `db:"kodepos" json:"kodepos"`

	AlamatDomisili  *string `db:"alamat_domisili" json:"alamat_domisili"`
	KabKotaDomID    *int `db:"kabkotaid_dom" json:"kabkota_dom_id"`
	KabKotaDom      *string `db:"kab_kota_dom" json:"kab_kota_dom"`
	ProvDomID       *int `db:"provid_dom" json:"prov_dom_id"`
	ProvinsiDom     *string `db:"provinsi_dom" json:"provinsi_dom"`
	KodePosDom      *string `db:"kodepos_dom" json:"kodepos_dom"`

	JnsPegawai *string `db:"jns_pegawai" json:"jenis_pegawai"`

	// ===== PENDIDIKAN UTAMA =====
	PddkID         *uuid.UUID `db:"pddkID" json:"pddk_id"`
	NamaJenjang    *string `db:"nama_jenjang" json:"nama_jenjang"`
	NamaInstitusi  *string `db:"nama_institusi" json:"nama_institusi"`
	TahunLulus     *int   `db:"tahun_lulus" json:"tahun_lulus"`

	// ===== HISTORY JSONB =====
	PendidikanHistory *json.RawMessage `db:"pendidikan_history" json:"pendidikan_history"`
	JabatanHistory    *json.RawMessage `db:"jabatan_history" json:"jabatan_history"`
	GradeHistory      *json.RawMessage `db:"grade_history" json:"grade_history"`
	KeluargaPegawai   *json.RawMessage `db:"keluarga_pegawai" json:"keluarga_pegawai"`
	KondarPegawai     *json.RawMessage `db:"kondar_pegawai" json:"kondar_pegawai"`

	// ===== JABATAN AKTIF =====
	DivisiJabatan *string     `db:"divisi_jabatan" json:"divisi_jabatan"`
	NamaJabatan   *string     `db:"nama_jabatan" json:"nama_jabatan"`
	TmtJab        *time.Time `db:"tmt_jab" json:"tmt_jab"`
	TmtJabAkhir   *time.Time `db:"tmt_jab_akhir" json:"tmt_jab_akhir"`
	Gapok         *float64   `db:"gapok" json:"gapok"`

	// ===== GROUP =====
	GroupNames *string        `db:"group_names" json:"group_names"`
	GroupIDsRaw pq.StringArray `db:"groupids" json:"-"`
	GroupIDs   []uuid.UUID  `db:"groupids" json:"group_ids"`

	// ===== GRADE AKTIF =====
	NamaGrade *string `db:"nama_grade" json:"nama_grade"`

	// ===== AUDIT =====
	CreatedBy   *string     `db:"created_by" json:"created_by"`
	CreatedDate *time.Time `db:"created_date" json:"created_date"`
	UpdatedBy   *string     `db:"updated_by" json:"updated_by"`
	UpdatedDate *time.Time `db:"updated_date" json:"updated_date"`
	StatusData  *string     `db:"status_data" json:"status_data"`
}