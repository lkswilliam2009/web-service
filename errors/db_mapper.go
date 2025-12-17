package errors

import (
	"github.com/lib/pq"
)

// Map PostgreSQL error → AppError
func FromDB(err error) *AppError {

	pqErr, ok := err.(*pq.Error)
	if !ok {
		// Bukan error PostgreSQL
		return Internal(err)
	}

	switch pqErr.Code {

	// UNIQUE VIOLATION
	case "23505":
		switch pqErr.Constraint {

		case "users_email_key":
			return Conflict(
				"DB_001",
				"Email already exists",
			)

		case "users_uname_key":
			return Conflict(
				"DB_002",
				"Username already exists",
			)

		default:
			return Conflict(
				"DB_003",
				"Duplicate data",
			)
		}

	// FOREIGN KEY VIOLATION
	case "23503":
		return Conflict(
			"DB_004",
			"Referenced data does not exist",
		)

	// NOT NULL VIOLATION
	case "23502":
		return BadRequest(
			"Required field is missing",
		)

	// CHECK VIOLATION
	case "23514":
		return BadRequest(
			"Invalid data format",
		)
	}

	// Default DB error
	return &AppError{
		HTTPCode: 500,
		Code:     "DB_999",
		Message:  "Database error",
		Err:      err,
	}
}