package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Port    int    `json:"port"`
	Env     string `json:"env"`
	Pepper  string `json:"pepper"`
	HMACKey string `json:"hmac_key"`

	Database PostgresConfig `json:"database"`
	Mailgun  MailgunConfig  `json:"mailgun"`
	Dropbox OAuthConfig `json:"dropbox"`
}

func DefaultConfig() Config {
	return Config{
		Port:     3000,
		Env:      "dev",
		Pepper:   "foobar",
		HMACKey:  "secret-hmac-key",
		Database: DefaultPostgresConfig(),
	}
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func LoadConfig(configReq bool) Config {
	// Open the config file
	f, err := os.Open(".config")
	if err != nil {
		if configReq {
			panic(err)
		}

		log.Println("Using the default config...")
		return DefaultConfig()
	}

	var c Config

	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully loaded .config")

	return c
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "192.168.56.101",
		Port:     5432,
		User:     "developer",
		Password: "1234qwer",
		Name:     "lenslockedbr_dev",
	}
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s "+
			"sslmode=disable", c.Host, c.Port,
			c.User, c.Name)
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", c.Host, c.Port,
		c.User, c.Password, c.Name)
}

type MailgunConfig struct {
	APIKey       string `json:"api_key"`
	PublicAPIKey string `json:"public_api_key"`
	Domain       string `json:"domain"`
}

type OAuthConfig struct {
	ID string `json:"id"`
	Secret string `json:"secret"`
	AuthURL string `json:"auth_url"`
	TokenURL string `json:"token_url"`
	RedirectURL string `json:"redirect_url"`
}


