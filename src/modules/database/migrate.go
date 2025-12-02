package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"slices"
	"strings"
)

//go:embed migrations/*.sql
var migrations embed.FS

func createMigrateTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		version		VARCHAR(36),
		applied_at	TIMESTAMPTZ	DEFAULT NOW()
	);`)
	if err != nil {
		return err
	}

	return nil
}

func getApplied(db *sql.DB) ([]string, error) {
	var applied []string = make([]string, 0)
	rows, err := db.Query("SELECT version FROM migrations;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}

		applied = append(applied, version)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return applied, nil
}

func (d *Database) migrate(db *sql.DB) error {
	rows, err := db.Query("SELECT version FROM migrations;")
	if err != nil {
		return err
	}
	defer rows.Close()

	var entries []fs.DirEntry
	entries, err = fs.ReadDir(migrations, "migrations")
	if err != nil {
		return err
	}

	err = createMigrateTable(db)
	if err != nil {
		return err
	}

	var applied []string = make([]string, 0)
	applied, err = getApplied(db)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		var version, _ = strings.CutSuffix(entry.Name(), ".sql")
		if slices.Contains(applied, version) {
			continue
		}

		fmt.Printf("[%s] Applying migration: %s\n", d.Name(), version)

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		var paths = fmt.Sprintf("migrations/%s", entry.Name())
		bytes, err := migrations.ReadFile(paths)
		if err != nil {
			tx.Rollback()
			return err
		}

		var sql = string(bytes)
		var queries = strings.Split(sql, ";")
		for _, raw := range queries {
			var query = strings.TrimSpace(raw)
			if query == "" {
				continue
			}

			_, err = tx.Exec(query)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		_, err = tx.Exec("INSERT INTO migrations (version) VALUES ($1);", version)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		fmt.Printf("[%s] Applied migration: %s", d.Name(), version)
	}

	return nil
}
