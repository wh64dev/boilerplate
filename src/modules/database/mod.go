package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"git.wh64.net/devproje/devproje-boilerplate/src/config"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Database struct {
	DB *sql.DB
}

func uri(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}

func intUri(key string, value int64) string {
	return fmt.Sprintf("%s=%d", key, value)
}

func nonrequired(key, value string) string {
	if value == "" {
		return ""
	}

	return fmt.Sprintf("%s=%s", key, value)
}

func Migration(db *sql.DB) error {
	var applied []string = make([]string, 0)
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		version		VARCHAR(36),
		applied_at	TIMESTAMPTZ	DEFAULT NOW()
	);`)
	if err != nil {
		return fmt.Errorf("error occurred when creating migration table: %v", err)
	}

	rows, err := db.Query("SELECT version FROM migrations;")
	if err != nil {
		return fmt.Errorf("failed to load applied migrations: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("failed to scan migration version: %v", err)
		}

		applied = append(applied, version)
	}

	if rows.Err(); err != nil {
		return fmt.Errorf("error during rows iteration: %v", err)
	}

	var entries []fs.DirEntry
	entries, err = fs.ReadDir(migrations, "migrations")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		// var version, _ = strings.CutSuffix(entry.Name(), ".sql")

	}

	return nil
}

func (db *Database) Name() string {
	return "database"
}

func (db *Database) Init() error {
	var dbconf = config.Get.Database
	var host = uri("host", dbconf.Host)
	var port = intUri("port", dbconf.Port)
	var dbname = uri("dbname", dbconf.Name)
	var user = nonrequired("user", dbconf.Username)
	var password = uri("password", dbconf.Password)
	var c = strings.Trim(fmt.Sprintf("sslmode=disable %s %s %s %s %s", host, port, dbname, user, password), " ")
	var d, err = sql.Open("postgres", c)
	d.SetMaxIdleConns(10)
	if err != nil {
		return err
	}
	defer d.Close()

	err = d.Ping()
	if err != nil {
		return err
	}

	db.DB = d
	return nil
}

func (db *Database) Destroy() error {
	if db.DB != nil {
		db.DB.Close()
	}

	db.DB = nil
	return nil
}

var DatabaseModule = &Database{}
