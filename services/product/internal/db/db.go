package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func StartDB() (*sql.DB, error) {
	dsn := "user=admin password=secret host=localhost port=5432 dbname=product_service sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed creating connection to DB: %s", err)
	}

	//defer db.Close()
	db.Ping()

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	return db, nil
}
