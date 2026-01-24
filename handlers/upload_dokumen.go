package handlers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func UploadDokumen(c *fiber.Ctx) error {
	// =====================
	// JWT
	// =====================
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrUnauthorized
	}
	claims := token.Claims.(jwt.MapClaims)
	_ = claims // nanti bisa dipakai buat ownership

	// =====================
	// FORM DATA
	// =====================
	pegawaiID := c.FormValue("pegawai_id")
	docType := strings.ToUpper(c.FormValue("doc_type"))

	if pegawaiID == "" || docType == "" {
		return fiber.NewError(400, "pegawai_id and doc_type are required")
	}

	// validate UUID
	if _, err := uuid.Parse(pegawaiID); err != nil {
		return fiber.NewError(400, "Invalid pegawai_id")
	}

	allowedDoc := map[string]bool{
		"KTP":  true,
		"NPWP": true,
		"BPJS": true,
	}
	if !allowedDoc[docType] {
		return fiber.NewError(400, "Invalid doc_type")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(400, "file is required")
	}

	// =====================
	// FILE VALIDATION
	// =====================
	if file.Size > 2*1024*1024 {
		return fiber.NewError(400, "Max file size 2MB")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".pdf": true,
	}
	if !allowedExt[ext] {
		return fiber.NewError(400, "Invalid file extension")
	}

	// =====================
	// CREATE FOLDER
	// =====================
	baseDir := filepath.Join("uploads", "pegawai", pegawaiID, docType)

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fiber.NewError(500, "Failed create directory")
	}

	// =====================
	// FILE NAME
	// =====================
	filename := strings.ToLower(docType) + ext
	fullPath := filepath.Join(baseDir, filename)

	// overwrite allowed
	if err := c.SaveFile(file, fullPath); err != nil {
		return fiber.NewError(500, "Failed save file")
	}

	return c.JSON(fiber.Map{
		"message": "Document uploaded successfully",
		"pegawai_id": pegawaiID,
		"doc_type": docType,
		"file": filename,
	})
}

func GetDokumen(c *fiber.Ctx) error {
	pegawaiIDParam := c.Params("pegawaiID")
	docType := strings.ToUpper(c.Params("docType"))
	filename := c.Params("filename")

	// ======================
	// VALIDASI DOCTYPE
	// ======================
	allowedDoc := map[string]bool{
		"KTP":  true,
		"NPWP": true,
		"BPJS": true,
	}
	if !allowedDoc[docType] {
		return fiber.ErrBadRequest
	}

	// ======================
	// JWT
	// ======================
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)

	role, _ := claims["role"].(string)
	jwtPegawaiID, _ := claims["pegawai_id"].(string)

	// ======================
	// OWNERSHIP CHECK
	// ======================
	isPrivileged := role == "831ae134-fe8c-4d75-a532-eaed3e6c1617" || role == "be242b62-c412-4229-b5ea-1906840d6eee"

	if !isPrivileged {
		if jwtPegawaiID == "" || jwtPegawaiID != pegawaiIDParam {
			return fiber.ErrForbidden
		}
	}

	// ======================
	// BUILD FILE PATH
	// ======================
	filePath := filepath.Join(
		"uploads",
		"pegawai",
		pegawaiIDParam,
		docType,
		filename,
	)

	// ======================
	// FILE EXIST CHECK
	// ======================
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fiber.ErrNotFound
	}

	// ======================
	// SEND FILE
	// ======================
	return c.SendFile(filePath)
}