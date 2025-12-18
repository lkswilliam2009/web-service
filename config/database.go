package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=data_induk_db port=5432 sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}
}
