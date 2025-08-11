package student

const (
	StudentIdKey     contextKey = "studentId"
	sqliteDriverName            = "sqlite3"

	// env variables
	dbPath      = "DB_PATH"
	databaseUrl = "DATABASE_URL"

	// errors
	errStudentNotFound = "student not found"
)
