package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     int64
	Name     string
	Username string
	Password string
}

type ConfigRaw struct {
	Host     string
	Port     int64
	Database *DatabaseConfig
}

var Get *ConfigRaw

func load(key string, defaults string) string {
	var ret = os.Getenv(key)
	if ret == "" {
		return defaults
	}

	return ret
}

func loadInt(key string, defaults int64) int64 {
	var ret = os.Getenv(key)
	if ret == "" {
		return defaults
	}

	var n, _ = strconv.ParseInt(ret, 10, 64)
	return n
}

func init() {
	var err = godotenv.Load()
	if err != nil {
		log.Fatalln("error loading .env file")
	}

	var host = load("HOST", "127.0.0.1")
	var port = loadInt("PORT", 3000)
	var dbHost = load("DB_HOST", "127.0.0.1")
	var dbPort = loadInt("DB_PORT", 5432)
	var dbName = load("DB_NAME", "boilerplate")
	var dbUsername = load("DB_USERNAME", "")
	var dbPassword = load("DB_PASSWORD", "")

	Get = &ConfigRaw{
		Host: host,
		Port: port,
		Database: &DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			Name:     dbName,
			Username: dbUsername,
			Password: dbPassword,
		},
	}
}
