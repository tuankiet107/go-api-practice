package config

import (
	"net"
	"net/url"
	"os"
)

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadDatabase() Database {
	return Database{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "go_api_practice"),
	}
}

func (c Database) DSN() string {
	dsn := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   c.Name,
	}

	query := dsn.Query()
	query.Set("sslmode", "disable")
	dsn.RawQuery = query.Encode()
	return dsn.String()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
