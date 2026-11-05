package database

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL, migrationsPath string) error {
	// Ensure isolated migrations table for content-service so it doesn't conflict with other services
	migrationDBURL := databaseURL
	if strings.Contains(migrationDBURL, "?") {
		if !strings.Contains(migrationDBURL, "x-migrations-table") {
			migrationDBURL += "&x-migrations-table=schema_migrations_content"
		}
	} else {
		migrationDBURL += "?x-migrations-table=schema_migrations_content"
	}

	parsedURL, _ := url.Parse(databaseURL)
	if parsedURL != nil && parsedURL.User != nil {
		log.Printf("Connecting for database migrations: %s@%s%s", parsedURL.User.Username(), parsedURL.Host, parsedURL.Path)
	}

	m, err := migrate.New("file://"+migrationsPath, migrationDBURL)
	if err != nil {
		return fmt.Errorf("migration/init: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("migrator source close error: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("migrator db close error: %v", dbErr)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Content database migration: no schema changes to apply.")
			return nil
		}
		return fmt.Errorf("migration/up: %w", err)
	}

	log.Println("Content database migrations executed successfully.")
	return nil
}
