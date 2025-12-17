package errors

import "github.com/lib/pq"

func DBError(err error) *AppError {
	if pqErr, ok := err.(*pq.Error); ok {

		if pqErr.Code == "23505" { // unique violation
			switch pqErr.Constraint {
			case "users_uname_key":
				return Conflict("DB_001", "Username already exists")
			case "users_email_key":
				return Conflict("DB_001", "Email already exists")
			default:
				return Conflict("DB_001", "Duplicate data")
			}
		}

		if pqErr.Code == "23503" {
			return BadRequest("Invalid reference data")
		}
	}

	return &AppError{
		HTTPCode: 500,
		Code:     "DB_999",
		Message:  "Database error",
		Err:      err,
	}
}
