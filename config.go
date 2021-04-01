package main

import "fmt"

// Base Configs
type Config struct {
	Port    int    `json:"port"`
	Env     string `json:"env"`
	Pepper  string `json:"pepper"`
	HMACKey string `json:"hmac_key"`
}

func (c Config) IsProduction() bool {
	return c.Env == "PRODUCTION"
}

func (c Config) IsDevelopment() bool {
	return c.Env == "DEVELOPMENT"
}

func DefaultConfig() Config {
	return Config{
		Port:    8005,
		Env:     "DEVELOPMENT",
		Pepper:  "secret-random-string",
		HMACKey: "secret-hmac-key",
	}
}

// Database Configs
type PostgresConfig struct {
	DBHost     string `json:"host"`
	DBPort     int    `json:"port"`
	DBUser     string `json:"user"`
	DBPassword string `json:"password"`
	DBName     string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.DBPassword == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.DBHost, c.DBPort, c.DBUser, c.DBName,
		)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "robert",
		DBPassword: "password",
		DBName:     "gallerio",
	}
}
