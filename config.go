package main

import (
	"fmt"
)

type PostgresConfig struct {
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Name string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s " +
                                   "sslmode=disable", c.Host, c.Port,
                                   c.User, c.Name)
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s " + 
                           "dbname=%s sslmode=disable", c.Host, c.Port,
                           c.User, c.Password, c.Name)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig {
		Host: "192.168.56.101",
		Port: 5432,
		User: "developer",
		Password: "1234qwer",
		Name: "lenslockedbr_dev",
	}
}

type Config struct {
	Port int `json:"port"`
	Env string `json:"env"`
	Pepper string `json:"pepper"`
	HMACKey string `json:"hmac_key"`
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func DefaultConfig() Config {
	return Config {
		Port: 3000,
		Env: "dev",
		Pepper: "foobar",
		HMACKey: "secret-hmac-key", 
	}
}
