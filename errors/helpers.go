package errors

import "net/http"

func BadRequest(msg string) *AppError {
	return &AppError{
		HTTPCode: http.StatusBadRequest,
		Code:     "REQ_001", 
		Message:  msg,
	}
}

func Unauthorized(msg string) *AppError {
	return &AppError{
		HTTPCode: http.StatusUnauthorized,
		Code:     "AUTH_001",
		Message:  msg,
	}
}

func Forbidden(msg string) *AppError {
	return &AppError{
		HTTPCode: http.StatusForbidden,
		Code:     "PERM_001",
		Message:  msg,
	}
}

func Conflict(code, msg string) *AppError {
	return &AppError{
		HTTPCode: http.StatusConflict,
		Code:     code, // contoh: DB_001
		Message:  msg,
	}
}

func Internal(err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusInternalServerError,
		Code:     "SYS_001",
		Message:  "Internal server error",
		Err:      err,
	}
}

func InvalidJSON(err error) *AppError {
	return &AppError{
		HTTPCode: http.StatusBadRequest,
		Code:     "REQ_001",
		Message:  "Invalid JSON payload",
		Err:      err,
	}
}