package main

import (
	"Unnispick/cmd/api"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"strings"
)

func main() {
	var migrationDir string
	var dbURL string
	var command string

	// Parse command line arguments
	flag.StringVar(&migrationDir, "path", "migrations", "Directory where migration files are stored")
	flag.StringVar(&dbURL, "db", os.Getenv("DATABASE_URL"), "Database connection string (or use DATABASE_URL env var)")
	flag.StringVar(&command, "command", "", "Command to run (migrate/api)")
	flag.Parse()

	if command == "" {
		log.Fatal("Command is required (migrate/api)")
	}

	switch strings.ToLower(command) {
	case "migrate":
		handleMigration(migrationDir, dbURL)
	case "api":
		api.StartAPI()
	default:
		log.Fatalf("Invalid command: %s", command)
	}
}

func handleMigration(migrationDir, dbURL string) {
	if dbURL == "" {
		log.Fatal("Database URL is required for migration")
	}

	// Get migration direction (up/down) from args
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Migration direction (up/down) is required")
	}
	direction := strings.ToLower(args[0])

	// Create migration instance
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationDir),
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	// Execute migration
	switch direction {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("No migrations to apply")
				return
			}
			log.Fatal(err)
		}
		log.Println("Successfully applied up migrations")

	case "down":
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("No migrations to apply")
				return
			}
			log.Fatal(err)
		}
		log.Println("Successfully applied down migrations")

	default:
		log.Fatalf("Invalid migration direction: %s", direction)
	}
}
