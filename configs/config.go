package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

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

// Mailgun Configs
type MailgunConfig struct {
	Domain       string `json:"domain"`
	APIKey       string `json:"api_key"`
	PublicAPIKey string `json:"pub_api_key"`
}

func DefaultMailgunConfig() MailgunConfig {
	return MailgunConfig{}
}

// Base Configs
type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACKey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
	Mailgun  MailgunConfig  `json:"mailgun"`
}

func (c Config) IsProduction() bool {
	return c.Env == "PRODUCTION"
}

func (c Config) IsDevelopment() bool {
	return c.Env == "DEVELOPMENT"
}

func DefaultConfig() Config {
	return Config{
		Port:     8000,
		Env:      "DEVELOPMENT",
		Pepper:   "secret-random-string",
		HMACKey:  "secret-hmac-key",
		Database: DefaultPostgresConfig(),
		Mailgun:  DefaultMailgunConfig(),
	}
}

// Load Config
func LoadConfig(configReq bool) Config {
	f, err := os.Open("configs/.config.json")
	if err != nil {
		if configReq {
			panic(err)
		}
		fmt.Println("Loaded default configs......")
		return DefaultConfig()
	}

	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded configs......")
	return cfg
}
