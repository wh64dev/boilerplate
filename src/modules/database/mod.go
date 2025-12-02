package database

import (
	"database/sql"
	"fmt"
	"strings"

	"git.wh64.net/devproje/devproje-boilerplate/src/config"
	_ "github.com/lib/pq"
)

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

	err = db.migrate(d)
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
