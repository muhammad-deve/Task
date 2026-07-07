package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	fmt.Println("Running migrations for main database...")

	pgConfig := pool.Config()
	connString := pgConfig.ConnString()
	connString += "?sslmode=disable"

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("error creating db, sql.Open: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating main DB driver: %w", err)
	}

	migrationsPath := getMigrationsFolderPath()
	migrationsURL := (&url.URL{Scheme: "file", Path: filepath.ToSlash(migrationsPath)}).String()

	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}

	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("error reading migration version: %w", err)
	}

	latestVersion, err := latestMigrationVersion(migrationsPath)
	if err != nil {
		return fmt.Errorf("error reading migrations folder: %w", err)
	}

	if !dirty && currentVersion > latestVersion {
		fmt.Printf("Main DB migrations already at version %d.\n", currentVersion)
		return nil
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %w", err)
	}

	fmt.Println("Main DB migrations applied successfully.")
	return nil
}

func getMigrationsFolderPath() string {
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}
	return filepath.Join(basePath, "migrations")
}

func latestMigrationVersion(path string) (uint, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	var latest uint64
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		versionText, _, ok := strings.Cut(name, "_")
		if !ok {
			continue
		}

		version, err := strconv.ParseUint(versionText, 10, 64)
		if err != nil {
			continue
		}

		if version > latest {
			latest = version
		}
	}

	return uint(latest), nil
}
