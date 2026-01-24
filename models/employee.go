package models

import (
	"database/sql"
)	

type Employee struct {
	PegawaiID     string         `json:"pegawai_id"`
	UserID        string         `json:"userid"`
	NIP           string         `json:"NIP"`
	NamaLengkap   string         `json:"nama_lengkap"`
	DivisiJabatan sql.NullString `json:"divisi_jabatan"`
	NamaJabatan   sql.NullString `json:"nama_jabatan"`
	NamaGrade     sql.NullString `json:"nama_grade"`
	JnsPegawai    sql.NullString `json:"jns_pegawai"`
	StatusData    int            `json:"status_data"`
}

type EmployeeFnDTO struct {
	PegawaiID      string
	UserID         string
	NIP            string
	NamaLengkap    string
	DivisiJabatan sql.NullString `json:"divisi_jabatan"`
	NamaJabatan   sql.NullString `json:"nama_jabatan"`
	NamaGrade     sql.NullString `json:"nama_grade"`
	JnsPegawai    sql.NullString `json:"jns_pegawai"`
	StatusData     int
	TotalData      int
	Message        string
}