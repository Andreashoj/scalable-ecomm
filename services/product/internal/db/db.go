package db

import (
	"andreasho/scalable-ecomm/pgk"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func StartDB() (*sqlx.DB, error) {
	isDev := os.Getenv("ENV") == "DEV"
	var dsn string
	if isDev {
		dsn = "postgres://admin:secret@product-postgres:5432/product_service?sslmode=disable"
	} else {
		dsn = ""
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed creating connection to DB: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed pinging DB: %v", err)
	}

	if isDev {
		if err = pgk.MigrationsRunner(dsn, "internal/db/migrations"); err != nil {
			return nil, fmt.Errorf("failed running migrations: %v", err)
		}
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	return sqlx.NewDb(db, "postgres"), nil
}
