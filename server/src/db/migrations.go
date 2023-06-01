package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

// MigrationsUp run migrations from directory with specified for the taken sql connection
func MigrationsUp(conn *sql.DB, pathToMigrations string) error {

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Println(err)
	}

	source, err := (&file.File{}).Open(pathToMigrations)
	if err != nil {
		log.Println("Open ", err)
	}

	m, err := migrate.NewWithInstance("file", source, "remote_storage", driver)
	if err != nil {
		log.Println("New With Instance", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Println("migrate Up", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Println("version ", err)
	}

	fmt.Printf("Current migration version: %v, Dirty: %v\n", version, dirty)
	return err
}
