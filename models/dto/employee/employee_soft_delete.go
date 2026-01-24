package dto

type EmployeeSoftDeleteRequest struct {
	PegawaiID string `json:"pegawai_id" validate:"required"`
}