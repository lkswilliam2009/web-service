package errors

type AppError struct {
	HTTPCode int    // HTTP status (400, 401, 500)
	Code     string // ERROR CODE: DB_001, AUTH_001
	Message  string // User-facing message
	Err      error  // Internal error (log only)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}