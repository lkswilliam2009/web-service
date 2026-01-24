package dto

type EmployeeDetailRequest struct {
	PegawaiID    string `json:"pegawaiID" validate:"required"`
	NIP          string `json:"NIP" validate:"required"`
}