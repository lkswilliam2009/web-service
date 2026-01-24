package config

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	DBx *sqlx.DB
)

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

	DBx = sqlx.NewDb(DB, "postgres")

	log.Println("Database connected (DB & DBx ready)")
}