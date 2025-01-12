package main

import (
	"errors"
	"github.com/Code-xCartel/noxus-api-svc/config"
	"github.com/Code-xCartel/noxus-api-svc/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func main() {
	dbDriver, dbErr := db.NewSQLStorage(config.Envs.PostgresConnectionString)
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	driver, err := postgres.WithInstance(dbDriver, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, mErr := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "postgres", driver)
	if mErr != nil {
		log.Fatal(mErr)
	}
	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Printf("migrate up failed: %v", err)
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	}
}
