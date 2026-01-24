package handlers

import (
	"log"
	"strconv"
	"math"
	"time"

	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"encoding/json"
	"github.com/google/uuid"
	_"github.com/lib/pq"

	"web-service/config"
	"web-service/models"
	"web-service/models/dto/employee"

	appErr "web-service/errors"
)

func EmployeeData(c *fiber.Ctx) error {
	// ===== JWT =====
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fiber.ErrUnauthorized
	}

	role, _ := claims["role"].(string)
	branch, _ := claims["branch"].(string)

	// ===== GROUPS =====
	var groupIDs []uuid.UUID

	if raw, ok := claims["groups"]; ok {
		bytes, _ := json.Marshal(raw)

		var groups []models.GroupArr
		if err := json.Unmarshal(bytes, &groups); err == nil {
			for _, g := range groups {
				if g.GroupID != uuid.Nil {
					groupIDs = append(groupIDs, g.GroupID)
				}
			}

		}
	}

	var branchParam interface{}

	if branch != "" {
		branchParam = branch
	} else {
		branchParam = nil
	}

	// ===== PAGINATION =====
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// ===== CALL FUNCTION =====
	rows, err := config.DB.Query(`
		SELECT * FROM fn_get_employee_data(
			$1::uuid,
			$2::uuid,
			$3,
			$4
		)
	`, role, branchParam, limit, offset)

	if err != nil {
		log.Println("ERROR DETAIL:", err)
		return appErr.FromDB(err)
	}
	defer rows.Close()

	var (
		data    []models.Employee
		total   int
		message string
	)

	for rows.Next() {
		var dto models.EmployeeFnDTO

		if err := rows.Scan(
			&dto.PegawaiID,
			&dto.UserID,
			&dto.NIP,
			&dto.NamaLengkap,
			&dto.DivisiJabatan,
			&dto.NamaJabatan,
			&dto.NamaGrade,
			&dto.JnsPegawai,
			&dto.StatusData,
			&dto.TotalData,
			&dto.Message,
		); err != nil {
			return appErr.FromDB(err)
		}

		total = dto.TotalData
		message = dto.Message

		data = append(data, models.Employee{
			PegawaiID:     dto.PegawaiID,
			UserID:        dto.UserID,
			NIP:           dto.NIP,
			NamaLengkap:   dto.NamaLengkap,
			DivisiJabatan: dto.DivisiJabatan,
			NamaJabatan:   dto.NamaJabatan,
			NamaGrade:     dto.NamaGrade,
			JnsPegawai:    dto.JnsPegawai,
			StatusData:    dto.StatusData,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(fiber.Map{
		"message": message,
		"data":    data,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"count":       len(data),
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

func EmployeeDataDetail(c *fiber.Ctx) error {
	var req dto.EmployeeDetailRequest

	if err := c.BodyParser(&req); err != nil {
		return appErr.InvalidJSON(err)
	}

	if req.PegawaiID == "" || req.NIP == "" {
		return appErr.BadRequest("pegawaiID and NIP required")
	}

	var result dto.EmployeeDetailResponse

	err := config.DBx.Get(
		&result,
		`SELECT * FROM fn_get_employee_data_detail($1::uuid, $2::text)`,
		req.PegawaiID,
		req.NIP,
	)

	if err != nil {
		log.Println("ERROR DETAIL:", err)
		if err == sql.ErrNoRows {
			return appErr.BadRequest("Employee not found")
		}
		return appErr.FromDB(err)
	}

	// ===== CONVERT groupids =====
	var groupUUIDs []uuid.UUID
	for _, g := range result.GroupIDsRaw {
		uid, err := uuid.Parse(g)
		if err == nil {
			groupUUIDs = append(groupUUIDs, uid)
		}
	}
	result.GroupIDs = groupUUIDs

	return c.JSON(fiber.Map{
		"message": "Success get employee detail",
		"data":    result,
	})
}

func EmployeeInsert(c *fiber.Ctx) error {
	var req dto.EmployeeInsertRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("ERROR DETAIL:", err)
		return appErr.InvalidJSON(err)
	}

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	createdByStr := claims["user_id"].(string)
	createdBy, _ := uuid.Parse(createdByStr)
	userID, _ := uuid.Parse(req.UserID)

	if req.UserID == "" || req.NIP == "" || req.NamaLengkap == "" || req.NoKTP == "" {
		return appErr.BadRequest("user_id, nip, nama_lengkap, no_ktp wajib diisi")
	}

	var res struct {
		PegawaiID   uuid.UUID `db:"pegawai_id" json:"pegawai_id"`
		UserID      uuid.UUID `db:"userid" json:"user_id"`
		NIP         string    `db:"nip" json:"nip"`
		NamaLengkap string    `db:"nama_lengkap" json:"nama_lengkap"`
		StatusData  int       `db:"status_data" json:"status_data"`
	}

	var tglLahir *time.Time
	if req.TglLahir != nil && *req.TglLahir != "" {
		t, err := time.Parse("2006-01-02", *req.TglLahir)
		if err != nil {
			return appErr.BadRequest("format born date must YYYY-MM-DD")
		}
		tglLahir = &t
	}

	var tinggi *int
	if req.Tinggi != nil && *req.Tinggi != "" {
		t, err := strconv.Atoi(*req.Tinggi)
		if err != nil {
			return appErr.BadRequest("tall must be integer")
		}
		tinggi = &t
	}

	var berat *int
	if req.Berat != nil && *req.Berat != "" {
		b, err := strconv.Atoi(*req.Berat)
		if err != nil {
			return appErr.BadRequest("weight must be integer")
		}
		berat = &b
	}

	var kabkotaid *int
	if req.KabKotaID != nil && *req.KabKotaID != "" {
		b, err := strconv.Atoi(*req.KabKotaID)
		if err != nil {
			return appErr.BadRequest("kabkotaid must be integer")
		}
		kabkotaid = &b
	}

	var kabkotaid_dom *int
	if req.KabKotaDomID != nil && *req.KabKotaDomID != "" {
		b, err := strconv.Atoi(*req.KabKotaDomID)
		if err != nil {
			return appErr.BadRequest("kabkotaid_dom must be integer")
		}
		kabkotaid_dom = &b
	}

	err := config.DBx.Get(&res, `
		SELECT * FROM fn_employee_insert(
		  $1::uuid,
		  $2::varchar,
		  $3::varchar,
		  $4::varchar,
		  $5::date,
		  $6::jns_kelamin,
		  $7::varchar,
		  $8::agama,
		  $9::warga_negara,
		  $10::status_kawin,
		  $11::gol_darah,
		  $12::int,
		  $13::int,
		  $14::varchar,
		  $15::varchar,
		  $16::varchar,
		  $17::varchar,
		  $18::varchar,
		  $19::varchar,
		  $20::varchar,
		  $21::varchar,
		  $22::jns_pegawai,
		  $23::text,
		  $24::int,
		  $25::varchar,
		  $26::text,
		  $27::int,
		  $28::varchar,
		  $29::varchar,
		  $30::varchar,
		  $31::varchar,
		  $32::uuid
		)
	`,
		userID,
		req.NIP,
		req.NamaLengkap,
		req.TempatLahir,
		tglLahir,
		req.JnsKelamin,
		req.Suku,
		req.Agama,
		req.WargaNegara,
		req.StatusKawin,
		req.GolDarah,
		tinggi,
		berat,
		req.Rambut,
		req.BentukMuka,
		req.WarnaKulit,
		req.CacatTubuh,
		req.Hobby,
		req.NoKTP,
		req.NoNPWP,
		req.NoBPJS,
		req.JnsPegawai,
		req.AlamatKTP,
		kabkotaid,
		req.KodePos,
		req.AlamatDomisili,
		kabkotaid_dom,
		req.KodePosDom,
		req.NoHP1,
		req.NoHP2,
		req.Email,
		createdBy,
	)

	if err != nil {
		log.Println("DB ERROR DETAIL:", err)
		return appErr.FromDB(err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Employee successfully created",
		"data": res,
	})
}

func EmployeeUpdate(c *fiber.Ctx) error {
	var body dto.EmployeeUpdateRequest

	if err := c.BodyParser(&body); err != nil {
		log.Println("ERROR DETAIL:", err)
		return appErr.InvalidJSON(err)
	}

	// ================= VALIDASI =================
	if body.PegawaiID == "" {
		return appErr.BadRequest("pegawai_id is required")
	}

	validGender := map[string]bool{
	  "L": true,
	  "P": true,
	}

	if body.JnsKelamin != nil && !validGender[*body.JnsKelamin] {
	  return appErr.BadRequest("Invalid gender")
	}

	var tglLahir *time.Time
	if body.TglLahir != nil && *body.TglLahir != "" {
		t, err := time.Parse("2006-01-02", *body.TglLahir)
		if err != nil {
			return appErr.BadRequest("format born date must YYYY-MM-DD")
		}
		tglLahir = &t
	}

	var tinggi *int
	if body.Tinggi != nil && *body.Tinggi != "" {
		t, err := strconv.Atoi(*body.Tinggi)
		if err != nil {
			return appErr.BadRequest("tall must be integer")
		}
		tinggi = &t
	}

	var berat *int
	if body.Berat != nil && *body.Berat != "" {
		b, err := strconv.Atoi(*body.Berat)
		if err != nil {
			return appErr.BadRequest("weight must be integer")
		}
		berat = &b
	}

	var kabkotaid *int
	if body.KabKotaID != nil && *body.KabKotaID != "" {
		b, err := strconv.Atoi(*body.KabKotaID)
		if err != nil {
			return appErr.BadRequest("kabkotaid must be integer")
		}
		kabkotaid = &b
	}

	var kabkotaid_dom *int
	if body.KabKotaDomID != nil && *body.KabKotaDomID != "" {
		b, err := strconv.Atoi(*body.KabKotaDomID)
		if err != nil {
			return appErr.BadRequest("kabkotaid_dom must be integer")
		}
		kabkotaid_dom = &b
	}

	// ================= JWT =================
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fiber.ErrUnauthorized
	}

	updatedByStr, ok := claims["user_id"].(string)
	if !ok || updatedByStr == "" {
		return fiber.ErrUnauthorized
	}

	updatedBy, err := uuid.Parse(updatedByStr)
	if err != nil {
		return appErr.BadRequest("invalid user_id")
	}

	// ================= PARSE ID =================
	pegawaiID, err := uuid.Parse(body.PegawaiID)
	if err != nil {
		return appErr.BadRequest("invalid pegawai_id")
	}

	userID, err := uuid.Parse(body.UserID)
	if err != nil {
		return appErr.BadRequest("invalid user_id")
	}

	// ================= RESPONSE =================
	type EmployeeUpdateResponse struct {
		PegawaiID   uuid.UUID `db:"pegawai_id" json:"pegawai_id"`
		UserID      uuid.UUID `db:"userid" json:"user_id"`
		NIP         string    `db:"nip" json:"nip"`
		NamaLengkap string    `db:"nama_lengkap" json:"nama_lengkap"`
		StatusData  int       `db:"status_data" json:"status_data"`
	}

	var res EmployeeUpdateResponse

	// ================= EXEC =================
	err = config.DBx.Get(&res, `
	  SELECT * FROM fn_employee_update(
	    $1::uuid,
		$2::uuid,
		$3::varchar,
		$4::varchar,
		$5::varchar,
		$6::date,
		$7::jns_kelamin,
		$8::varchar,
		$9::agama,
		$10::warga_negara,
		$11::status_kawin,
		$12::gol_darah,
		$13::int,
		$14::int,
		$15::varchar,
		$16::varchar,
		$17::varchar,
		$18::varchar,
		$19::varchar,
		$20::varchar,
		$21::varchar, 
		$22::varchar,
		$23::jns_pegawai,
		$24::text, 
		$25::int, 
		$26::varchar,
		$27::text,
		$28::int,
		$29::varchar,
		$30::varchar,
		$31::varchar,
		$32::varchar,
		$33::uuid
	  )
	`,
	  pegawaiID,
	  userID,
	  body.NIP,
	  body.NamaLengkap,
	  body.TempatLahir,
	  tglLahir,
	  body.JnsKelamin,
	  body.Suku,
	  body.Agama,
	  body.WargaNegara,
	  body.StatusKawin,
	  body.GolDarah,
	  tinggi,
	  berat,
	  body.Rambut,
	  body.BentukMuka,
	  body.WarnaKulit,
	  body.CacatTubuh,
	  body.Hobby,
	  body.NoKTP,
	  body.NoNPWP,
	  body.NoBPJS,
	  body.JnsPegawai, 
	  body.AlamatKTP,
	  kabkotaid,
	  body.KodePos,
	  body.AlamatDomisili,
	  kabkotaid_dom,
	  body.KodePosDom,
	  body.NoHP1,
	  body.NoHP2,
	  body.Email,
	  updatedBy,
	)

	if err != nil {
		log.Println("DB ERROR DETAIL:", err)
		return appErr.FromDB(err)
	}

	return c.JSON(fiber.Map{
		"message": "Employee successfully updated",
		"data":    res,
	})
}

func EmployeeSoftDelete(c *fiber.Ctx) error {
	var body dto.EmployeeSoftDeleteRequest

	if err := c.BodyParser(&body); err != nil {
		return appErr.InvalidJSON(err)
	}

	if body.PegawaiID == "" {
		return appErr.BadRequest("pegawai_id required")
	}

	// ===== JWT =====
	token, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fiber.ErrUnauthorized
	}

	updatedBy, ok := claims["user_id"].(string)
	if !ok || updatedBy == "" {
		return fiber.ErrUnauthorized
	}

	// ===== UPDATE SOFT DELETE =====
	result, err := config.DB.Exec(`
		UPDATE master_pegawai
		SET
			status_data = 5,
			updated_by = $1,
			updated_at = NOW()
		WHERE "pegawaiID" = $2
	`, updatedBy, body.PegawaiID)

	if err != nil {
		return appErr.FromDB(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return appErr.FromDB(err)
	}

	if rows == 0 {
		return appErr.BadRequest("Employee not found or already inactive")
	}

	// ===== RESPONSE =====
	return c.JSON(fiber.Map{
		"message": "Employee successfully deactivated",
		"pegawai_id": body.PegawaiID,
	})
}