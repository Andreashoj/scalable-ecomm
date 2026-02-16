package db

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestDB(t *testing.T) *sql.DB {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_DB":       "product_test_db",
		},
		Cmd: []string{
			"postgres",
			"-c", "listen_addresses=*",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Errorf("failed to start container: %v", err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://postgres:pass@%s:%s/product_test_db?sslmode=disable", host, port.Port())

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to start db: %v", err)
	}

	if err = db.Ping(); err != nil {
		t.Fatalf("failed trying to ping db: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatalf("failed to create driver: %v", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	migrationPath := filepath.Join(filepath.Dir(filename), "../db/migrations")
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationPath, "postgres", driver)
	if err != nil {
		t.Fatalf("failed to create migrator: %v", err)
	}

	if err = m.Up(); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	t.Cleanup(func() {
		//container.Terminate(ctx)
		t.Logf("Container still running at %s:%s", host, port.Port())
	})

	return db
}
